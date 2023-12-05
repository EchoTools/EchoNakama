const errMissingPayload: nkruntime.Error = {
  message: 'no payload provided.',
  code: nkruntime.Codes.NOT_FOUND
};
const errBadInput: nkruntime.Error = {
  message: 'input contained invalid data.',
  code: nkruntime.Codes.INVALID_ARGUMENT
};
const errMissingId: nkruntime.Error = {
  message: 'no "id" provided.',
  code: nkruntime.Codes.NOT_FOUND
};

const errInternal = (message: string): nkruntime.Error => ({
  message: 'errInternal: ' + message,
  code: nkruntime.Codes.INTERNAL
} as nkruntime.Error);

export {
  errBadInput,
  errInternal,
  errMissingId,
  errMissingPayload,
}