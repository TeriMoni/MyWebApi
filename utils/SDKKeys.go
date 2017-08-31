package utils

func GetAttackNames() []string {
	attackNames := []string{"xposed", "substrate", "libinject", "adbi", "dagger", "roboguice", "butterknife", "dexposed", "poison", "ptrace"}
	return attackNames
}

func GetDangerPackages() []string {

	dangerPackages := []string{"com.doubee.mtkmaster", "com.example.myxposed", "de.robv.android.xposed.installer", "org.imei.mtk65xx", "com.pwdgame.imeisi",
		"cn.mm.gk", "com.soft.apk008", "com.mockgps.outside.ui", "locationcheater", "com.songg.version", "com.songzi.mmodel", "com.soft.apk008v",
		"com.sollyu.xposed.hook.model", "com.unidevel.devicemod", "com.wireless.macchangerqybm", "com.jiaofamily.android.mac", "com.lemonsqueeze.fakewificonnection", "xposed",
		"xprivacy", "org.tracetool.hackconnectivityservice"}
	return dangerPackages
}
