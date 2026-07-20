package main

import "github.com/Baba01hacker666/Go2APK/ui"

func init() {
	ui.RegisterEvent("btn_3_onclick", on3)
	ui.RegisterEvent("btn_clear_onclick", onClear)
	ui.RegisterEvent("btn_mul_onclick", onMul)
	ui.RegisterEvent("btn_1_onclick", on1)
	ui.RegisterEvent("btn_add_onclick", onAdd)
	ui.RegisterEvent("btn_0_onclick", on0)
	ui.RegisterEvent("btn_dot_onclick", onDot)
	ui.RegisterEvent("btn_9_onclick", on9)
	ui.RegisterEvent("btn_del_onclick", onDel)
	ui.RegisterEvent("btn_eq_onclick", onEq)
	ui.RegisterEvent("btn_rp_onclick", onRP)
	ui.RegisterEvent("btn_7_onclick", on7)
	ui.RegisterEvent("btn_8_onclick", on8)
	ui.RegisterEvent("btn_4_onclick", on4)
	ui.RegisterEvent("btn_6_onclick", on6)
	ui.RegisterEvent("btn_sub_onclick", onSub)
	ui.RegisterEvent("btn_2_onclick", on2)
	ui.RegisterEvent("btn_lp_onclick", onLP)
	ui.RegisterEvent("btn_div_onclick", onDiv)
	ui.RegisterEvent("btn_5_onclick", on5)
}
