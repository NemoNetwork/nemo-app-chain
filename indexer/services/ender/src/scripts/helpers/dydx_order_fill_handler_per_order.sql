CREATE OR REPLACE FUNCTION dydx_order_fill_handler_per_order(
    field text, block_height int, block_time timestamp, event_data jsonb, event_index int, transaction_index int,
    transaction_hash text, fill_liquidity text, fill_type text, usdc_asset_id text, order_canceled_status text) RETURNS jsonb AS $$
/**
  Parameters:
    - field: the field storing the order to process.
    - block_height: the height of the block being processing.
    - block_time: the time of the block being processed.
    - event_data: The 'data' field of the IndexerTendermintEvent (https://github.com/nemo-network/v4-chain/blob/9ed26bd/proto/nemo-network/indexer/indexer_manager/event.proto#L25)
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
    - event_index: The 'event_index' of the IndexerTendermintEvent.
    - transaction_index: The transaction_index of the IndexerTendermintEvent after the conversion that takes into
        account the block_event (https://github.com/nemo-network/v4-chain/blob/9ed26bd/indexer/services/ender/src/lib/helper.ts#L41)
    - transaction_hash: The transaction hash corresponding to this event from the IndexerTendermintBlock 'tx_hashes'.
    - fill_liquidity: The liquidity for the fill record.
    - fill_type: The type for the fill record.
    - usdc_asset_id: The USDC asset id.
    - order_canceled_status: Status of order cancelation
  Returns: JSON object containing fields:
    - order: The updated order in order-model format (https://github.com/nemo-network/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/order-model.ts).
    - fill: The updated fill in fill-model format (https://github.com/nemo-network/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/fill-model.ts).
    - perpetual_market: The perpetual market for the order in perpetual-market-model format (https://github.com/nemo-network/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/perpetual-market-model.ts).
    - perpetual_position: The updated perpetual position in perpetual-position-model format (https://github.com/nemo-network/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/perpetual-position-model.ts).

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    order_ jsonb;
    maker_order jsonb;
    clob_pair_id bigint;
    subaccount_uuid uuid;
    perpetual_market_record perpetual_markets%ROWTYPE;
    market_record markets%ROWTYPE;
    order_record orders%ROWTYPE;
    fill_record fills%ROWTYPE;
    perpetual_position_record perpetual_positions%ROWTYPE;
    asset_record assets%ROWTYPE;
    order_uuid uuid;
    order_side text;
    order_size numeric;
    order_price numeric;
    order_client_metadata bigint;
    fee numeric;
    fill_amount numeric;
    total_filled numeric;
    maker_price numeric;
    event_id bytea;
BEGIN
    order_ = event_data->field;
    maker_order = event_data->'makerOrder';
    clob_pair_id = jsonb_extract_path(order_, 'orderId', 'clobPairId')::bigint;
    perpetual_market_record = dydx_get_perpetual_market_for_clob_pair(clob_pair_id);
    market_record = dydx_get_market_for_id(perpetual_market_record."marketId");

    BEGIN
        SELECT * INTO STRICT asset_record FROM assets WHERE "id" = usdc_asset_id;
    EXCEPTION
        WHEN NO_DATA_FOUND THEN
            RAISE EXCEPTION 'Unable to find asset with id %', usdc_asset_id;
    END;

    /**
      Calculate sizes, prices, and fill amounts.

      TODO(IND-238): Extract out calculation of quantums and subticks to their own SQL functions.
    */
    order_size = dydx_trim_scale(dydx_from_jsonlib_long(order_->'quantums') *
                                 power(10, perpetual_market_record."atomicResolution")::numeric);
    order_price = dydx_trim_scale(dydx_from_jsonlib_long(order_->'subticks') *
                                  power(10, perpetual_market_record."quantumConversionExponent" +
                                                     asset_record."atomicResolution" -
                                                     perpetual_market_record."atomicResolution")::numeric);
    fill_amount = dydx_trim_scale(dydx_from_jsonlib_long(event_data->'fillAmount') *
                                  power(10, perpetual_market_record."atomicResolution")::numeric);
    maker_price = dydx_trim_scale(dydx_from_jsonlib_long(maker_order->'subticks') *
                                  power(10, perpetual_market_record."quantumConversionExponent" +
                                                     asset_record."atomicResolution" -
                                                     perpetual_market_record."atomicResolution")::numeric);
    total_filled = dydx_trim_scale(dydx_get_total_filled(fill_liquidity, event_data) *
                                   power(10, perpetual_market_record."atomicResolution")::numeric);
    fee = dydx_trim_scale(dydx_get_fee(fill_liquidity, event_data) *
                          power(10, asset_record."atomicResolution")::numeric);

    order_uuid = dydx_uuid_from_order_id(order_->'orderId');
    subaccount_uuid = dydx_uuid_from_subaccount_id(jsonb_extract_path(order_, 'orderId', 'subaccountId'));
    order_side = dydx_from_protocol_order_side(order_->'side');
    order_client_metadata = (order_->'clientMetadata')::bigint;

    /** Upsert the order, populating the order_record fields with what will be in the database. */
    SELECT * INTO order_record FROM orders WHERE "id" = order_uuid;
    order_record."side" = order_side;
    order_record."size" = order_size;
    order_record."price" = order_price;
    order_record."timeInForce" = dydx_from_protocol_time_in_force(order_->'timeInForce');
    order_record."reduceOnly" = (order_->>'reduceOnly')::boolean;
    order_record."orderFlags" = jsonb_extract_path(order_, 'orderId', 'orderFlags')::bigint;
    order_record."goodTilBlock" = (order_->'goodTilBlock')::bigint;
    order_record."goodTilBlockTime" = to_timestamp((order_->'goodTilBlockTime')::double precision);
    order_record."clientMetadata" = order_client_metadata;
    order_record."updatedAt" = block_time;
    order_record."updatedAtHeight" = block_height;

    IF FOUND THEN
        order_record."totalFilled" = total_filled;
        order_record."status" = dydx_get_order_status(total_filled, order_record.size, order_canceled_status, order_record."orderFlags", order_record."timeInForce");

        UPDATE orders
        SET
            "side" = order_record."side",
            "size" = order_record."size",
            "totalFilled" = order_record."totalFilled",
            "price" = order_record."price",
            "status" = order_record."status",
            "orderFlags" = order_record."orderFlags",
            "goodTilBlock" = order_record."goodTilBlock",
            "goodTilBlockTime" = order_record."goodTilBlockTime",
            "timeInForce" = order_record."timeInForce",
            "reduceOnly" = order_record."reduceOnly",
            "clientMetadata" = order_record."clientMetadata",
            "updatedAt" = order_record."updatedAt",
            "updatedAtHeight" = order_record."updatedAtHeight"
        WHERE id = order_uuid;
    ELSE
        order_record."id" = order_uuid;
        order_record."subaccountId" = subaccount_uuid;
        order_record."clientId" = jsonb_extract_path_text(order_, 'orderId', 'clientId')::bigint;
        order_record."clobPairId" = clob_pair_id;
        order_record."side" = order_side;
        order_record."type" = 'LIMIT'; /* TODO: Add additional order types once we support */

        order_record."totalFilled" = fill_amount;
        order_record."status" = dydx_get_order_status(fill_amount, order_size, order_canceled_status, order_record."orderFlags", order_record."timeInForce");
        order_record."createdAtHeight" = block_height;
        INSERT INTO orders
            ("id", "subaccountId", "clientId", "clobPairId", "side", "size", "totalFilled", "price", "type",
            "status", "timeInForce", "reduceOnly", "orderFlags", "goodTilBlock", "goodTilBlockTime", "createdAtHeight",
            "clientMetadata", "triggerPrice", "updatedAt", "updatedAtHeight")
        VALUES (order_record.*);
    END IF;

    /* Insert the associated fill record for this order_fill event. */
    event_id = dydx_event_id_from_parts(
        block_height, transaction_index, event_index);
    INSERT INTO fills
        ("id", "subaccountId", "side", "liquidity", "type", "clobPairId", "orderId", "size", "price", "quoteAmount",
         "eventId", "transactionHash", "createdAt", "createdAtHeight", "clientMetadata", "fee")
    VALUES (dydx_uuid_from_fill_event_parts(event_id, fill_liquidity),
            subaccount_uuid,
            order_side,
            fill_liquidity,
            fill_type,
            clob_pair_id,
            order_uuid,
            fill_amount,
            maker_price,
            dydx_trim_scale(fill_amount * maker_price),
            event_id,
            transaction_hash,
            block_time,
            block_height,
            order_client_metadata,
            fee)
    RETURNING * INTO fill_record;

    /* Upsert the perpetual_position record for this order_fill event. */
    perpetual_position_record = dydx_update_perpetual_position_aggregate_fields(
            subaccount_uuid,
            perpetual_market_record."id",
            order_side,
            fill_amount,
            maker_price);

    RETURN jsonb_build_object(
            'order',
            dydx_to_jsonb(order_record),
            'fill',
            dydx_to_jsonb(fill_record),
            'perpetual_market',
            dydx_to_jsonb(perpetual_market_record),
            'market',
            dydx_to_jsonb(market_record),
            'perpetual_position',
            dydx_to_jsonb(perpetual_position_record)
        );
END;
$$ LANGUAGE plpgsql;
