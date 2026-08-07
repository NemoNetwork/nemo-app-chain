import { bigIntToBytes } from '@nemo-network-indexer/v4-proto-parser';

export function intToSerializedInt(int: number): Uint8Array {
  return bigIntToBytes(BigInt(int));
}
