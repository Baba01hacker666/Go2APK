package main

import "github.com/Baba01hacker666/Go2APK/android"

func init() {
	android.RegisterEvent("btn_add_onclick", onAdd)
	android.RegisterEvent("btn_0_onclick", on0)
	android.RegisterEvent("btn_mul_onclick", onMul)
	android.RegisterEvent("btn_1_onclick", on1)
	android.RegisterEvent("btn_3_onclick", on3)
	android.RegisterEvent("btn_div_onclick", onDiv)
	android.RegisterEvent("btn_9_onclick", on9)
	android.RegisterEvent("btn_4_onclick", on4)
	android.RegisterEvent("btn_6_onclick", on6)
	android.RegisterEvent("btn_2_onclick", on2)
	android.RegisterEvent("btn_dot_onclick", onDot)
	android.RegisterEvent("btn_del_onclick", onDel)
	android.RegisterEvent("btn_rp_onclick", onRP)
	android.RegisterEvent("btn_7_onclick", on7)
	android.RegisterEvent("btn_sub_onclick", onSub)
	android.RegisterEvent("btn_eq_onclick", onEq)
	android.RegisterEvent("btn_clear_onclick", onClear)
	android.RegisterEvent("btn_lp_onclick", onLP)
	android.RegisterEvent("btn_8_onclick", on8)
	android.RegisterEvent("btn_5_onclick", on5)
}
