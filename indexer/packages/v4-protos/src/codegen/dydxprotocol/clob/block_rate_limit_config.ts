import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Defines the block rate limits for CLOB specific operations. */

export interface BlockRateLimitConfiguration {
  maxShortTermOrdersPerMarketPerNBlocks: MaxPerNBlocksRateLimit[];
  /**
   * How many stateful order attempts (successful and failed) are allowed for a
   * subaccount per N blocks. Note that the rate limits are applied
   * in an AND fashion such that an order must pass all rate limit
   * configurations.
   * 
   * Specifying 0 values disables this rate limit.
   */

  maxStatefulOrdersPerNBlocks: MaxPerNBlocksRateLimit[];
  maxShortTermOrderCancellationsPerMarketPerNBlocks: MaxPerNBlocksRateLimit[];
}
/** Defines the block rate limits for CLOB specific operations. */

export interface BlockRateLimitConfigurationSDKType {
  max_short_term_orders_per_market_per_n_blocks: MaxPerNBlocksRateLimitSDKType[];
  /**
   * How many stateful order attempts (successful and failed) are allowed for a
   * subaccount per N blocks. Note that the rate limits are applied
   * in an AND fashion such that an order must pass all rate limit
   * configurations.
   * 
   * Specifying 0 values disables this rate limit.
   */

  max_stateful_orders_per_n_blocks: MaxPerNBlocksRateLimitSDKType[];
  max_short_term_order_cancellations_per_market_per_n_blocks: MaxPerNBlocksRateLimitSDKType[];
}
/** Defines a rate limit over a specific number of blocks. */

export interface MaxPerNBlocksRateLimit {
  /**
   * How many blocks the rate limit is over.
   * Specifying 0 is invalid.
   */
  numBlocks: number;
  /**
   * What the limit is for `num_blocks`.
   * Specifying 0 is invalid.
   */

  limit: number;
}
/** Defines a rate limit over a specific number of blocks. */

export interface MaxPerNBlocksRateLimitSDKType {
  /**
   * How many blocks the rate limit is over.
   * Specifying 0 is invalid.
   */
  num_blocks: number;
  /**
   * What the limit is for `num_blocks`.
   * Specifying 0 is invalid.
   */

  limit: number;
}

function createBaseBlockRateLimitConfiguration(): BlockRateLimitConfiguration {
  return {
    maxShortTermOrdersPerMarketPerNBlocks: [],
    maxStatefulOrdersPerNBlocks: [],
    maxShortTermOrderCancellationsPerMarketPerNBlocks: []
  };
}

export const BlockRateLimitConfiguration = {
  encode(message: BlockRateLimitConfiguration, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.maxShortTermOrdersPerMarketPerNBlocks) {
      MaxPerNBlocksRateLimit.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.maxStatefulOrdersPerNBlocks) {
      MaxPerNBlocksRateLimit.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    for (const v of message.maxShortTermOrderCancellationsPerMarketPerNBlocks) {
      MaxPerNBlocksRateLimit.encode(v!, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BlockRateLimitConfiguration {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockRateLimitConfiguration();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.maxShortTermOrdersPerMarketPerNBlocks.push(MaxPerNBlocksRateLimit.decode(reader, reader.uint32()));
          break;

        case 2:
          message.maxStatefulOrdersPerNBlocks.push(MaxPerNBlocksRateLimit.decode(reader, reader.uint32()));
          break;

        case 3:
          message.maxShortTermOrderCancellationsPerMarketPerNBlocks.push(MaxPerNBlocksRateLimit.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<BlockRateLimitConfiguration>): BlockRateLimitConfiguration {
    const message = createBaseBlockRateLimitConfiguration();
    message.maxShortTermOrdersPerMarketPerNBlocks = object.maxShortTermOrdersPerMarketPerNBlocks?.map(e => MaxPerNBlocksRateLimit.fromPartial(e)) || [];
    message.maxStatefulOrdersPerNBlocks = object.maxStatefulOrdersPerNBlocks?.map(e => MaxPerNBlocksRateLimit.fromPartial(e)) || [];
    message.maxShortTermOrderCancellationsPerMarketPerNBlocks = object.maxShortTermOrderCancellationsPerMarketPerNBlocks?.map(e => MaxPerNBlocksRateLimit.fromPartial(e)) || [];
    return message;
  }

};

function createBaseMaxPerNBlocksRateLimit(): MaxPerNBlocksRateLimit {
  return {
    numBlocks: 0,
    limit: 0
  };
}

export const MaxPerNBlocksRateLimit = {
  encode(message: MaxPerNBlocksRateLimit, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.numBlocks !== 0) {
      writer.uint32(8).uint32(message.numBlocks);
    }

    if (message.limit !== 0) {
      writer.uint32(16).uint32(message.limit);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MaxPerNBlocksRateLimit {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMaxPerNBlocksRateLimit();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.numBlocks = reader.uint32();
          break;

        case 2:
          message.limit = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MaxPerNBlocksRateLimit>): MaxPerNBlocksRateLimit {
    const message = createBaseMaxPerNBlocksRateLimit();
    message.numBlocks = object.numBlocks ?? 0;
    message.limit = object.limit ?? 0;
    return message;
  }

};