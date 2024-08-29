import express from 'express';

import AddressesController from './addresses-controller';
import AssetPositionsController from './asset-positions-controller';
import CandlesController from './candles-controller';
import ComplianceController from './compliance-controller';
import ComplianceV2Controller from './compliance-v2-controller';
import FillsController from './fills-controller';
import HeightController from './height-controller';
import HistoricalBlockTradingRewardController from './historical-block-trading-rewards-controller';
import HistoricalFundingController from './historical-funding-controller';
import PnlticksController from './historical-pnl-controller';
import HistoricalTradingRewardController from './historical-trading-reward-aggregations-controller';
import OrderbooksController from './orderbook-controller';
import OrdersController from './orders-controller';
import PerpetualMarketController from './perpetual-markets-controller';
import PerpetualPositionsController from './perpetual-positions-controller';
import SparklinesController from './sparklines-controller';
import TimeController from './time-controller';
import TradesController from './trades-controller';
import TransfersController from './transfers-controller';

// Keep routers in alphabetical order

const router: express.Router = express.Router();
router.use('/addresses', AddressesController);
router.use('/assetPositions', AssetPositionsController);
router.use('/candles', CandlesController);
router.use('/fills', FillsController);
router.use('/height', HeightController);
router.use('/historicalBlockTradingRewards', HistoricalBlockTradingRewardController);
router.use('/historicalFunding', HistoricalFundingController);
router.use('/historical-pnl', PnlticksController);
router.use('/historicalTradingRewardAggregations', HistoricalTradingRewardController);
router.use('/orders', OrdersController);
router.use('/orderbooks', OrderbooksController);
router.use('/perpetualMarkets', PerpetualMarketController);
router.use('/perpetualPositions', PerpetualPositionsController);
router.use('/sparklines', SparklinesController);
router.use('/time', TimeController);
router.use('/trades', TradesController);
router.use('/transfers', TransfersController);
router.use('/screen', ComplianceController);
router.use('/compliance', ComplianceV2Controller);

export default router;
