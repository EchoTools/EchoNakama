import {
	setAccountRpc,
	getAccountRpc,
	setConfigRpc,
	getConfigRpc,
	setDocumentRpc,
	getDocumentRpc,
	setAccessControlListRpc,
	getAccessControlListRpc,
	setChannelInfoRpc,
	getChannelInfoRpc,
	setLoginSettingsRpc,
	getLoginSettingsRpc,
  getDeviceLinkCodeRpc,
  discordLinkDeviceRpc
} from './rpc';

// Initialize the server module
let InitModule: nkruntime.InitModule =
  function (ctx: nkruntime.Context,
    logger: nkruntime.Logger,
    nk: nkruntime.Nakama,
    initializer: nkruntime.Initializer) {

    initializer.registerRpc('echorelay/setAccount', setAccountRpc);
    initializer.registerRpc('echorelay/getAccount', getAccountRpc);
    initializer.registerRpc('echorelay/setConfig', setConfigRpc);
    initializer.registerRpc('echorelay/getConfig', getConfigRpc);
    initializer.registerRpc('echorelay/setDocument', setDocumentRpc);
    initializer.registerRpc('echorelay/getDocument', getDocumentRpc);
    initializer.registerRpc('echorelay/setAccessControlList', setAccessControlListRpc);
    initializer.registerRpc('echorelay/getAccessControlList', getAccessControlListRpc);
    initializer.registerRpc('echorelay/setChannelInfo', setChannelInfoRpc);
    initializer.registerRpc('echorelay/getChannelInfo', getChannelInfoRpc);
    initializer.registerRpc('echorelay/setLoginSettings', setLoginSettingsRpc);
    initializer.registerRpc('echorelay/getLoginSettings', getLoginSettingsRpc);
    initializer.registerRpc('echorelay/getDeviceLinkCode', getDeviceLinkCodeRpc);
    initializer.registerRpc('discordLinkDevice', discordLinkDeviceRpc);
  }


// Reference InitModule to avoid it getting removed on build
!InitModule && InitModule.bind(null);
