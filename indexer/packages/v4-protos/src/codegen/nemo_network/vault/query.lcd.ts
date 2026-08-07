import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryParamsRequest, QueryParamsResponseSDKType, QueryVaultRequest, QueryVaultResponseSDKType, QueryAllVaultsRequest, QueryAllVaultsResponseSDKType, QueryMegavaultTotalSharesRequest, QueryMegavaultTotalSharesResponseSDKType, QueryMegavaultOwnerSharesRequest, QueryMegavaultOwnerSharesResponseSDKType, QueryMegavaultAllOwnerSharesRequest, QueryMegavaultAllOwnerSharesResponseSDKType, QueryVaultParamsRequest, QueryVaultParamsResponseSDKType, QueryMegavaultFeeStateRequest, QueryMegavaultFeeStateResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.params = this.params.bind(this);
    this.vault = this.vault.bind(this);
    this.allVaults = this.allVaults.bind(this);
    this.megavaultTotalShares = this.megavaultTotalShares.bind(this);
    this.megavaultOwnerShares = this.megavaultOwnerShares.bind(this);
    this.megavaultAllOwnerShares = this.megavaultAllOwnerShares.bind(this);
    this.vaultParams = this.vaultParams.bind(this);
    this.megavaultFeeState = this.megavaultFeeState.bind(this);
  }
  /* Queries the Params. */


  async params(_params: QueryParamsRequest = {}): Promise<QueryParamsResponseSDKType> {
    const endpoint = `nemo_network/vault/params`;
    return await this.req.get<QueryParamsResponseSDKType>(endpoint);
  }
  /* Queries a Vault by type and number. */


  async vault(params: QueryVaultRequest): Promise<QueryVaultResponseSDKType> {
    const endpoint = `nemo_network/vault/vault/${params.type}/${params.number}`;
    return await this.req.get<QueryVaultResponseSDKType>(endpoint);
  }
  /* Queries all vaults. */


  async allVaults(params: QueryAllVaultsRequest = {
    pagination: undefined
  }): Promise<QueryAllVaultsResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `nemo_network/vault/vault`;
    return await this.req.get<QueryAllVaultsResponseSDKType>(endpoint, options);
  }
  /* Queries total shares of megavault. */


  async megavaultTotalShares(_params: QueryMegavaultTotalSharesRequest = {}): Promise<QueryMegavaultTotalSharesResponseSDKType> {
    const endpoint = `nemo_network/vault/megavault/total_shares`;
    return await this.req.get<QueryMegavaultTotalSharesResponseSDKType>(endpoint);
  }
  /* Queries owner shares of megavault. */


  async megavaultOwnerShares(params: QueryMegavaultOwnerSharesRequest): Promise<QueryMegavaultOwnerSharesResponseSDKType> {
    const endpoint = `nemo_network/vault/megavault/owner_shares/${params.address}`;
    return await this.req.get<QueryMegavaultOwnerSharesResponseSDKType>(endpoint);
  }
  /* Queries all owner shares of megavault. */


  async megavaultAllOwnerShares(params: QueryMegavaultAllOwnerSharesRequest = {
    pagination: undefined
  }): Promise<QueryMegavaultAllOwnerSharesResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `nemo_network/vault/megavault/all_owner_shares`;
    return await this.req.get<QueryMegavaultAllOwnerSharesResponseSDKType>(endpoint, options);
  }
  /* Queries vault params of a vault. */


  async vaultParams(params: QueryVaultParamsRequest): Promise<QueryVaultParamsResponseSDKType> {
    const endpoint = `nemo_network/vault/params/${params.type}/${params.number}`;
    return await this.req.get<QueryVaultParamsResponseSDKType>(endpoint);
  }
  /* Queries the megavault fee accounting state.
   Note: fork-local; has no upstream equivalent. */


  async megavaultFeeState(_params: QueryMegavaultFeeStateRequest = {}): Promise<QueryMegavaultFeeStateResponseSDKType> {
    const endpoint = `nemo_network/vault/megavault/fee_state`;
    return await this.req.get<QueryMegavaultFeeStateResponseSDKType>(endpoint);
  }

}