/* ------- USER COMPLAINT TYPES ------- */

type IsoString = string;

export interface UserComplaintCreateObject {
  id: string,
  walletAddress: string,
  message: string,
  email?: string,
  createdAt?: IsoString,
}

export enum UserComplaintColumns {
  id = 'id',
  walletAddress = 'walletAddress',
  message = 'message',
  email = 'email',
  createdAt = 'createdAt',
}
