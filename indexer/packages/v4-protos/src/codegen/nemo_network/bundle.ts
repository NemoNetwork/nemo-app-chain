import * as _12 from "./accountplus/accountplus";
import * as _13 from "./accountplus/genesis";
import * as _14 from "./affiliates/affiliates";
import * as _15 from "./affiliates/genesis";
import * as _16 from "./affiliates/query";
import * as _17 from "./affiliates/tx";
import * as _18 from "./assets/asset";
import * as _19 from "./assets/genesis";
import * as _20 from "./assets/query";
import * as _21 from "./assets/tx";
import * as _22 from "./blocktime/blocktime";
import * as _23 from "./blocktime/genesis";
import * as _24 from "./blocktime/params";
import * as _25 from "./blocktime/query";
import * as _26 from "./blocktime/tx";
import * as _27 from "./bridge/bridge_event_info";
import * as _28 from "./bridge/bridge_event";
import * as _29 from "./bridge/genesis";
import * as _30 from "./bridge/params";
import * as _31 from "./bridge/query";
import * as _32 from "./bridge/tx";
import * as _33 from "./clob/block_rate_limit_config";
import * as _34 from "./clob/clob_pair";
import * as _35 from "./clob/equity_tier_limit_config";
import * as _36 from "./clob/genesis";
import * as _37 from "./clob/liquidations_config";
import * as _38 from "./clob/liquidations";
import * as _39 from "./clob/matches";
import * as _40 from "./clob/mev";
import * as _41 from "./clob/operation";
import * as _42 from "./clob/order_removals";
import * as _43 from "./clob/order";
import * as _44 from "./clob/process_proposer_matches_events";
import * as _45 from "./clob/query";
import * as _46 from "./clob/tx";
import * as _47 from "./daemons/bridge/bridge";
import * as _48 from "./daemons/liquidation/liquidation";
import * as _49 from "./daemons/pricefeed/price_feed";
import * as _50 from "./delaymsg/block_message_ids";
import * as _51 from "./delaymsg/delayed_message";
import * as _52 from "./delaymsg/genesis";
import * as _53 from "./delaymsg/query";
import * as _54 from "./delaymsg/tx";
import * as _55 from "./epochs/epoch_info";
import * as _56 from "./epochs/genesis";
import * as _57 from "./epochs/query";
import * as _58 from "./feetiers/genesis";
import * as _59 from "./feetiers/params";
import * as _60 from "./feetiers/query";
import * as _61 from "./feetiers/tx";
import * as _62 from "./govplus/genesis";
import * as _63 from "./govplus/query";
import * as _64 from "./govplus/tx";
import * as _65 from "./indexer/events/events";
import * as _66 from "./indexer/indexer_manager/event";
import * as _67 from "./indexer/off_chain_updates/off_chain_updates";
import * as _68 from "./indexer/protocol/v1/clob";
import * as _69 from "./indexer/protocol/v1/perpetual";
import * as _70 from "./indexer/protocol/v1/subaccount";
import * as _71 from "./indexer/protocol/v1/vault";
import * as _72 from "./indexer/redis/redis_order";
import * as _73 from "./indexer/shared/removal_reason";
import * as _74 from "./indexer/socks/messages";
import * as _75 from "./listing/genesis";
import * as _76 from "./listing/params";
import * as _77 from "./listing/query";
import * as _78 from "./listing/tx";
import * as _79 from "./perpetuals/genesis";
import * as _80 from "./perpetuals/params";
import * as _81 from "./perpetuals/perpetual";
import * as _82 from "./perpetuals/query";
import * as _83 from "./perpetuals/tx";
import * as _84 from "./prices/genesis";
import * as _85 from "./prices/market_param";
import * as _86 from "./prices/market_price";
import * as _87 from "./prices/query";
import * as _88 from "./prices/tx";
import * as _89 from "./ratelimit/capacity";
import * as _90 from "./ratelimit/genesis";
import * as _91 from "./ratelimit/limit_params";
import * as _92 from "./ratelimit/pending_send_packet";
import * as _93 from "./ratelimit/query";
import * as _94 from "./ratelimit/tx";
import * as _95 from "./revshare/genesis";
import * as _96 from "./revshare/params";
import * as _97 from "./revshare/query";
import * as _98 from "./revshare/revshare";
import * as _99 from "./revshare/tx";
import * as _100 from "./rewards/genesis";
import * as _101 from "./rewards/params";
import * as _102 from "./rewards/query";
import * as _103 from "./rewards/reward_share";
import * as _104 from "./rewards/tx";
import * as _105 from "./sending/genesis";
import * as _106 from "./sending/query";
import * as _107 from "./sending/transfer";
import * as _108 from "./sending/tx";
import * as _109 from "./stats/genesis";
import * as _110 from "./stats/params";
import * as _111 from "./stats/query";
import * as _112 from "./stats/stats";
import * as _113 from "./stats/tx";
import * as _114 from "./subaccounts/asset_position";
import * as _115 from "./subaccounts/genesis";
import * as _116 from "./subaccounts/perpetual_position";
import * as _117 from "./subaccounts/query";
import * as _118 from "./subaccounts/streaming";
import * as _119 from "./subaccounts/subaccount";
import * as _120 from "./vault/genesis";
import * as _121 from "./vault/params";
import * as _122 from "./vault/query";
import * as _123 from "./vault/share";
import * as _124 from "./vault/tx";
import * as _125 from "./vault/vault";
import * as _126 from "./vest/genesis";
import * as _127 from "./vest/query";
import * as _128 from "./vest/tx";
import * as _129 from "./vest/vest_entry";
import * as _130 from "./assets/query.lcd";
import * as _131 from "./blocktime/query.lcd";
import * as _132 from "./bridge/query.lcd";
import * as _133 from "./clob/query.lcd";
import * as _134 from "./delaymsg/query.lcd";
import * as _135 from "./epochs/query.lcd";
import * as _136 from "./feetiers/query.lcd";
import * as _137 from "./listing/query.lcd";
import * as _138 from "./perpetuals/query.lcd";
import * as _139 from "./prices/query.lcd";
import * as _140 from "./ratelimit/query.lcd";
import * as _141 from "./revshare/query.lcd";
import * as _142 from "./rewards/query.lcd";
import * as _143 from "./stats/query.lcd";
import * as _144 from "./subaccounts/query.lcd";
import * as _145 from "./vault/query.lcd";
import * as _146 from "./vest/query.lcd";
import * as _147 from "./affiliates/query.rpc.Query";
import * as _148 from "./assets/query.rpc.Query";
import * as _149 from "./blocktime/query.rpc.Query";
import * as _150 from "./bridge/query.rpc.Query";
import * as _151 from "./clob/query.rpc.Query";
import * as _152 from "./delaymsg/query.rpc.Query";
import * as _153 from "./epochs/query.rpc.Query";
import * as _154 from "./feetiers/query.rpc.Query";
import * as _155 from "./govplus/query.rpc.Query";
import * as _156 from "./listing/query.rpc.Query";
import * as _157 from "./perpetuals/query.rpc.Query";
import * as _158 from "./prices/query.rpc.Query";
import * as _159 from "./ratelimit/query.rpc.Query";
import * as _160 from "./revshare/query.rpc.Query";
import * as _161 from "./rewards/query.rpc.Query";
import * as _162 from "./sending/query.rpc.Query";
import * as _163 from "./stats/query.rpc.Query";
import * as _164 from "./subaccounts/query.rpc.Query";
import * as _165 from "./vault/query.rpc.Query";
import * as _166 from "./vest/query.rpc.Query";
import * as _167 from "./affiliates/tx.rpc.msg";
import * as _168 from "./blocktime/tx.rpc.msg";
import * as _169 from "./bridge/tx.rpc.msg";
import * as _170 from "./clob/tx.rpc.msg";
import * as _171 from "./delaymsg/tx.rpc.msg";
import * as _172 from "./feetiers/tx.rpc.msg";
import * as _173 from "./govplus/tx.rpc.msg";
import * as _174 from "./listing/tx.rpc.msg";
import * as _175 from "./perpetuals/tx.rpc.msg";
import * as _176 from "./prices/tx.rpc.msg";
import * as _177 from "./ratelimit/tx.rpc.msg";
import * as _178 from "./revshare/tx.rpc.msg";
import * as _179 from "./rewards/tx.rpc.msg";
import * as _180 from "./sending/tx.rpc.msg";
import * as _181 from "./stats/tx.rpc.msg";
import * as _182 from "./vault/tx.rpc.msg";
import * as _183 from "./vest/tx.rpc.msg";
import * as _184 from "./lcd";
import * as _185 from "./rpc.query";
import * as _186 from "./rpc.tx";
export namespace nemo_network {
  export const accountplus = { ..._12,
    ..._13
  };
  export const affiliates = { ..._14,
    ..._15,
    ..._16,
    ..._17,
    ..._147,
    ..._167
  };
  export const assets = { ..._18,
    ..._19,
    ..._20,
    ..._21,
    ..._130,
    ..._148
  };
  export const blocktime = { ..._22,
    ..._23,
    ..._24,
    ..._25,
    ..._26,
    ..._131,
    ..._149,
    ..._168
  };
  export const bridge = { ..._27,
    ..._28,
    ..._29,
    ..._30,
    ..._31,
    ..._32,
    ..._132,
    ..._150,
    ..._169
  };
  export const clob = { ..._33,
    ..._34,
    ..._35,
    ..._36,
    ..._37,
    ..._38,
    ..._39,
    ..._40,
    ..._41,
    ..._42,
    ..._43,
    ..._44,
    ..._45,
    ..._46,
    ..._133,
    ..._151,
    ..._170
  };
  export namespace daemons {
    export const bridge = { ..._47
    };
    export const liquidation = { ..._48
    };
    export const pricefeed = { ..._49
    };
  }
  export const delaymsg = { ..._50,
    ..._51,
    ..._52,
    ..._53,
    ..._54,
    ..._134,
    ..._152,
    ..._171
  };
  export const epochs = { ..._55,
    ..._56,
    ..._57,
    ..._135,
    ..._153
  };
  export const feetiers = { ..._58,
    ..._59,
    ..._60,
    ..._61,
    ..._136,
    ..._154,
    ..._172
  };
  export const govplus = { ..._62,
    ..._63,
    ..._64,
    ..._155,
    ..._173
  };
  export namespace indexer {
    export const events = { ..._65
    };
    export const indexer_manager = { ..._66
    };
    export const off_chain_updates = { ..._67
    };
    export namespace protocol {
      export const v1 = { ..._68,
        ..._69,
        ..._70,
        ..._71
      };
    }
    export const redis = { ..._72
    };
    export const shared = { ..._73
    };
    export const socks = { ..._74
    };
  }
  export const listing = { ..._75,
    ..._76,
    ..._77,
    ..._78,
    ..._137,
    ..._156,
    ..._174
  };
  export const perpetuals = { ..._79,
    ..._80,
    ..._81,
    ..._82,
    ..._83,
    ..._138,
    ..._157,
    ..._175
  };
  export const prices = { ..._84,
    ..._85,
    ..._86,
    ..._87,
    ..._88,
    ..._139,
    ..._158,
    ..._176
  };
  export const ratelimit = { ..._89,
    ..._90,
    ..._91,
    ..._92,
    ..._93,
    ..._94,
    ..._140,
    ..._159,
    ..._177
  };
  export const revshare = { ..._95,
    ..._96,
    ..._97,
    ..._98,
    ..._99,
    ..._141,
    ..._160,
    ..._178
  };
  export const rewards = { ..._100,
    ..._101,
    ..._102,
    ..._103,
    ..._104,
    ..._142,
    ..._161,
    ..._179
  };
  export const sending = { ..._105,
    ..._106,
    ..._107,
    ..._108,
    ..._162,
    ..._180
  };
  export const stats = { ..._109,
    ..._110,
    ..._111,
    ..._112,
    ..._113,
    ..._143,
    ..._163,
    ..._181
  };
  export const subaccounts = { ..._114,
    ..._115,
    ..._116,
    ..._117,
    ..._118,
    ..._119,
    ..._144,
    ..._164
  };
  export const vault = { ..._120,
    ..._121,
    ..._122,
    ..._123,
    ..._124,
    ..._125,
    ..._145,
    ..._165,
    ..._182
  };
  export const vest = { ..._126,
    ..._127,
    ..._128,
    ..._129,
    ..._146,
    ..._166,
    ..._183
  };
  export const ClientFactory = { ..._184,
    ..._185,
    ..._186
  };
}