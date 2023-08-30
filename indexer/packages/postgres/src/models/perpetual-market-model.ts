import path from 'path';

import { Model } from 'objection';

import {
  IntegerPattern,
  NonNegativeNumericPattern,
  NumericPattern,
} from '../lib/validators';
import {
  PerpetualMarketStatus,
} from '../types';

export default class PerpetualMarketModel extends Model {
  static get tableName() {
    return 'perpetual_markets';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    perpetualPosition: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'perpetual-position-model'),
      join: {
        from: 'perpetual_markets.id',
        to: 'perpetual_positions.perpetualId',
      },
    },
    market: {
      relation: Model.HasOneRelation,
      modelClass: path.join(__dirname, 'market-model'),
      join: {
        from: 'perpetual_markets.marketId',
        to: 'markets.id',
      },
    },
    liquidityTiers: {
      relation: Model.HasOneRelation,
      modelClass: path.join(__dirname, 'liquidity-tiers-model'),
      join: {
        from: 'perpetual_markets.liquidityTierId',
        to: 'liquidity_tiers.id',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'clobPairId',
        'ticker',
        'marketId',
        'status',
        'baseAsset',
        'quoteAsset',
        'lastPrice',
        'priceChange24H',
        'volume24H',
        'trades24H',
        'nextFundingRate',
        'basePositionSize',
        'incrementalPositionSize',
        'maxPositionSize',
        'openInterest',
        'quantumConversionExponent',
        'atomicResolution',
        'subticksPerTick',
        'minOrderBaseQuantums',
        'stepBaseQuantums',
        'liquidityTierId',
      ],
      properties: {
        id: { type: 'string', pattern: IntegerPattern },
        clobPairId: { type: 'string', pattern: IntegerPattern },
        ticker: { type: 'string' },
        marketId: { type: 'integer' },
        status: { type: 'string', enum: [...Object.values(PerpetualMarketStatus)] },
        baseAsset: { type: 'string' },
        quoteAsset: { type: 'string' },
        lastPrice: { type: 'string', pattern: NonNegativeNumericPattern },
        priceChange24H: { type: 'string', pattern: NumericPattern },
        volume24H: { type: 'string', pattern: NonNegativeNumericPattern },
        trades24H: { type: 'integer' },
        nextFundingRate: { type: 'string', pattern: NumericPattern },
        basePositionSize: { type: 'string', pattern: NonNegativeNumericPattern },
        incrementalPositionSize: { type: 'string', pattern: NonNegativeNumericPattern },
        maxPositionSize: { type: 'string', pattern: NonNegativeNumericPattern },
        openInterest: { type: 'string', pattern: NumericPattern },
        quantumConversionExponent: { type: 'integer' },
        atomicResolution: { type: 'integer' },
        subticksPerTick: { type: 'integer' },
        minOrderBaseQuantums: { type: 'integer' },
        stepBaseQuantums: { type: 'integer' },
        liquidityTierId: { type: 'integer' },
      },
    };
  }

  /**
   * A mapping from column name to JSON conversion expected.
   * See getSqlConversionForDydxModelTypes for valid conversions.
   *
   * TODO(IND-239): Ensure that jsonSchema() / sqlToJsonConversions() / model fields match.
   */
  static get sqlToJsonConversions() {
    return {
      id: 'string',
      clobPairId: 'string',
      ticker: 'string',
      marketId: 'integer',
      status: 'string',
      baseAsset: 'string',
      quoteAsset: 'string',
      lastPrice: 'string',
      priceChange24H: 'string',
      volume24H: 'string',
      trades24H: 'integer',
      nextFundingRate: 'string',
      basePositionSize: 'string',
      incrementalPositionSize: 'string',
      maxPositionSize: 'string',
      openInterest: 'string',
      quantumConversionExponent: 'integer',
      atomicResolution: 'integer',
      subticksPerTick: 'integer',
      minOrderBaseQuantums: 'integer',
      stepBaseQuantums: 'integer',
      liquidityTierId: 'integer',
    };
  }

  id!: string;

  clobPairId!: string;

  ticker!: string;

  marketId!: number;

  status!: PerpetualMarketStatus;

  baseAsset!: string;

  quoteAsset!: string;

  lastPrice!: string;

  priceChange24H!: string;

  volume24H!: string;

  trades24H!: number;

  nextFundingRate!: string;

  basePositionSize!: string;

  incrementalPositionSize!: string;

  maxPositionSize!: string;

  openInterest!: string;

  quantumConversionExponent!: number;

  atomicResolution!: number;

  subticksPerTick!: number;

  minOrderBaseQuantums!: number;

  stepBaseQuantums!: number;

  liquidityTierId!: number;
}
