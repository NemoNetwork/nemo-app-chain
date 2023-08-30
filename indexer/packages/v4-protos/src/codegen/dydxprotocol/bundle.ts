import * as _5 from "./assets/asset";
import * as _6 from "./assets/genesis";
import * as _7 from "./assets/query";
import * as _8 from "./assets/tx";
import * as _9 from "./blocktime/blocktime";
import * as _10 from "./blocktime/genesis";
import * as _11 from "./blocktime/params";
import * as _12 from "./blocktime/query";
import * as _13 from "./blocktime/tx";
import * as _14 from "./bridge/bridge_event_info";
import * as _15 from "./bridge/bridge_event";
import * as _16 from "./bridge/genesis";
import * as _17 from "./bridge/params";
import * as _18 from "./bridge/query";
import * as _19 from "./bridge/tx";
import * as _20 from "./clob/block_rate_limit_config";
import * as _21 from "./clob/clob_pair";
import * as _22 from "./clob/equity_tier_limit_config";
import * as _23 from "./clob/genesis";
import * as _24 from "./clob/liquidations_config";
import * as _25 from "./clob/liquidations";
import * as _26 from "./clob/matches";
import * as _27 from "./clob/mev";
import * as _28 from "./clob/operation";
import * as _29 from "./clob/order_removals";
import * as _30 from "./clob/order";
import * as _31 from "./clob/process_proposer_matches_events";
import * as _32 from "./clob/query";
import * as _33 from "./clob/tx";
import * as _34 from "./daemons/bridge/bridge";
import * as _35 from "./daemons/liquidation/liquidation";
import * as _36 from "./daemons/pricefeed/price_feed";
import * as _37 from "./delaymsg/block_message_ids";
import * as _38 from "./delaymsg/delayed_message";
import * as _39 from "./delaymsg/genesis";
import * as _40 from "./delaymsg/query";
import * as _41 from "./delaymsg/tx";
import * as _42 from "./epochs/epoch_info";
import * as _43 from "./epochs/genesis";
import * as _44 from "./epochs/query";
import * as _45 from "./feetiers/genesis";
import * as _46 from "./feetiers/params";
import * as _47 from "./feetiers/query";
import * as _48 from "./feetiers/tx";
import * as _49 from "./indexer/events/events";
import * as _50 from "./indexer/indexer_manager/event";
import * as _51 from "./indexer/off_chain_updates/off_chain_updates";
import * as _52 from "./indexer/protocol/v1/clob";
import * as _53 from "./indexer/protocol/v1/subaccount";
import * as _54 from "./indexer/redis/redis_order";
import * as _55 from "./indexer/shared/removal_reason";
import * as _56 from "./indexer/socks/messages";
import * as _57 from "./perpetuals/genesis";
import * as _58 from "./perpetuals/params";
import * as _59 from "./perpetuals/perpetual";
import * as _60 from "./perpetuals/query";
import * as _61 from "./perpetuals/tx";
import * as _62 from "./prices/genesis";
import * as _63 from "./prices/market_param";
import * as _64 from "./prices/market_price";
import * as _65 from "./prices/query";
import * as _66 from "./prices/tx";
import * as _67 from "./rewards/genesis";
import * as _68 from "./rewards/params";
import * as _69 from "./rewards/query";
import * as _70 from "./rewards/reward_share";
import * as _71 from "./rewards/tx";
import * as _72 from "./sending/genesis";
import * as _73 from "./sending/query";
import * as _74 from "./sending/transfer";
import * as _75 from "./sending/tx";
import * as _76 from "./stats/genesis";
import * as _77 from "./stats/params";
import * as _78 from "./stats/query";
import * as _79 from "./stats/stats";
import * as _80 from "./stats/tx";
import * as _81 from "./subaccounts/asset_position";
import * as _82 from "./subaccounts/genesis";
import * as _83 from "./subaccounts/perpetual_position";
import * as _84 from "./subaccounts/query";
import * as _85 from "./subaccounts/subaccount";
import * as _86 from "./vest/genesis";
import * as _87 from "./vest/query";
import * as _88 from "./vest/tx";
import * as _89 from "./vest/vest_entry";
import * as _96 from "./blocktime/query.lcd";
import * as _97 from "./bridge/query.lcd";
import * as _98 from "./clob/query.lcd";
import * as _99 from "./delaymsg/query.lcd";
import * as _100 from "./epochs/query.lcd";
import * as _101 from "./feetiers/query.lcd";
import * as _102 from "./perpetuals/query.lcd";
import * as _103 from "./prices/query.lcd";
import * as _104 from "./rewards/query.lcd";
import * as _105 from "./stats/query.lcd";
import * as _106 from "./subaccounts/query.lcd";
import * as _107 from "./vest/query.lcd";
import * as _108 from "./assets/query.rpc.Query";
import * as _109 from "./blocktime/query.rpc.Query";
import * as _110 from "./bridge/query.rpc.Query";
import * as _111 from "./clob/query.rpc.Query";
import * as _112 from "./delaymsg/query.rpc.Query";
import * as _113 from "./epochs/query.rpc.Query";
import * as _114 from "./feetiers/query.rpc.Query";
import * as _115 from "./perpetuals/query.rpc.Query";
import * as _116 from "./prices/query.rpc.Query";
import * as _117 from "./rewards/query.rpc.Query";
import * as _118 from "./sending/query.rpc.Query";
import * as _119 from "./stats/query.rpc.Query";
import * as _120 from "./subaccounts/query.rpc.Query";
import * as _121 from "./vest/query.rpc.Query";
import * as _122 from "./blocktime/tx.rpc.msg";
import * as _123 from "./bridge/tx.rpc.msg";
import * as _124 from "./clob/tx.rpc.msg";
import * as _125 from "./delaymsg/tx.rpc.msg";
import * as _126 from "./feetiers/tx.rpc.msg";
import * as _127 from "./perpetuals/tx.rpc.msg";
import * as _128 from "./prices/tx.rpc.msg";
import * as _129 from "./rewards/tx.rpc.msg";
import * as _130 from "./sending/tx.rpc.msg";
import * as _131 from "./stats/tx.rpc.msg";
import * as _132 from "./vest/tx.rpc.msg";
import * as _133 from "./lcd";
import * as _134 from "./rpc.query";
import * as _135 from "./rpc.tx";
export namespace dydxprotocol {
  export const assets = { ..._5,
    ..._6,
    ..._7,
    ..._8,
    ..._108
  };
  export const blocktime = { ..._9,
    ..._10,
    ..._11,
    ..._12,
    ..._13,
    ..._96,
    ..._109,
    ..._122
  };
  export const bridge = { ..._14,
    ..._15,
    ..._16,
    ..._17,
    ..._18,
    ..._19,
    ..._97,
    ..._110,
    ..._123
  };
  export const clob = { ..._20,
    ..._21,
    ..._22,
    ..._23,
    ..._24,
    ..._25,
    ..._26,
    ..._27,
    ..._28,
    ..._29,
    ..._30,
    ..._31,
    ..._32,
    ..._33,
    ..._98,
    ..._111,
    ..._124
  };
  export namespace daemons {
    export const bridge = { ..._34
    };
    export const liquidation = { ..._35
    };
    export const pricefeed = { ..._36
    };
  }
  export const delaymsg = { ..._37,
    ..._38,
    ..._39,
    ..._40,
    ..._41,
    ..._99,
    ..._112,
    ..._125
  };
  export const epochs = { ..._42,
    ..._43,
    ..._44,
    ..._100,
    ..._113
  };
  export const feetiers = { ..._45,
    ..._46,
    ..._47,
    ..._48,
    ..._101,
    ..._114,
    ..._126
  };
  export namespace indexer {
    export const events = { ..._49
    };
    export const indexer_manager = { ..._50
    };
    export const off_chain_updates = { ..._51
    };
    export namespace protocol {
      export const v1 = { ..._52,
        ..._53
      };
    }
    export const redis = { ..._54
    };
    export const shared = { ..._55
    };
    export const socks = { ..._56
    };
  }
  export const perpetuals = { ..._57,
    ..._58,
    ..._59,
    ..._60,
    ..._61,
    ..._102,
    ..._115,
    ..._127
  };
  export const prices = { ..._62,
    ..._63,
    ..._64,
    ..._65,
    ..._66,
    ..._103,
    ..._116,
    ..._128
  };
  export const rewards = { ..._67,
    ..._68,
    ..._69,
    ..._70,
    ..._71,
    ..._104,
    ..._117,
    ..._129
  };
  export const sending = { ..._72,
    ..._73,
    ..._74,
    ..._75,
    ..._118,
    ..._130
  };
  export const stats = { ..._76,
    ..._77,
    ..._78,
    ..._79,
    ..._80,
    ..._105,
    ..._119,
    ..._131
  };
  export const subaccounts = { ..._81,
    ..._82,
    ..._83,
    ..._84,
    ..._85,
    ..._106,
    ..._120
  };
  export const vest = { ..._86,
    ..._87,
    ..._88,
    ..._89,
    ..._107,
    ..._121,
    ..._132
  };
  export const ClientFactory = { ..._133,
    ..._134,
    ..._135
  };
}