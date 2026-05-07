/* ------- USER COMPLAINT TYPES ------- */

type IsoString = string;

export interface UserComplaintCreateObject {
  id: string,
  message: string,
  email?: string,
  createdAt?: IsoString,
}

export enum UserComplaintColumns {
  id = 'id',
  message = 'message',
  email = 'email',
  createdAt = 'createdAt',
}
