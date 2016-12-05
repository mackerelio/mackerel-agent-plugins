package mpmulticore

import "testing"

func TestParseProcStats(t *testing.T) {
	stab := `cpu  25308301 0 19470191 35582590482 432542 4227 1237778 2053417 0 0
cpu0 14786890 0 4570343 2196397805 364717 4072 1209169 526216 1 2
cpu1 2031368 0 2810375 2220960345 44625 21 6306 215873 0 0
cpu2 1136965 0 2003305 2223401399 16913 34 2490 150020 0 0
cpu3 887262 0 1351478 2225031553 1842 34 1984 106188 0 0
cpu4 698240 0 1022997 2225790144 688 5 1726 92769 0 0
cpu5 633207 0 855102 2226094460 512 7 1473 94372 0 0
cpu6 560396 0 772765 2226295609 465 4 1375 86695 0 0
cpu7 550697 0 717794 2226396684 314 4 1358 87411 0 0
cpu8 516476 0 684738 2226480673 388 4 1358 83375 0 0
cpu9 539329 0 690613 2226463131 329 4 1553 85872 0 0
cpu10 511186 0 681326 2226525771 323 5 1532 78527 0 0
cpu11 523958 0 689858 2226494064 280 6 1644 85290 0 0
cpu12 503041 0 696754 2226525183 313 4 1647 81350 0 0
cpu13 491136 0 652495 2226560016 268 9 1368 91025 0 0
cpu14 463343 0 635362 2226607769 291 3 1351 89751 0 0
cpu15 474802 0 634877 2226565867 266 6 1434 98675 0 0
intr 10189662361 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 514520202 67947609 233900 0 275291 698846453 217070356 243231 0 279760 640805807 126705906 256807 0 308045 324533905 96494233 262605 0 307572 247100682 86989223 263747 0 329175 231787116 86942679 262067 0 345493 211239797 83799100 257524 0 335479 201229182 80301598 258051 0 337002 192338035 78887243 261288 0 324064 194020233 77276543 264627 0 313805 184420804 77598965 268260 0 306992 188987565 76213860 270250 0 315551 181737220 76019767 268350 0 302552 191458500 76433789 266589 0 299170 188638355 77097745 266078 0 292540 201812993 77509758 263478 0 289432 494 300 67773621 285 2583576515 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 8676303806
btime 1444899785
processes 30982630
procs_running 2
procs_blocked 0
softirq 8788848490 0 2700996173 1094252449 1510759801 0 0 1 1761055606 1241616 1720542844`

	stat, _ := parseProcStat(stab)
	if len(stat) != 16 {
		t.Errorf("parseProcStat: size should be 16, but '%d'", len(stat))
	}
	if *stat["cpu0"].User != 14786890 {
		t.Errorf("parseProcStat: user should be 14786890, but '%f'", *stat["cpu0"].User)
	}
	if *stat["cpu0"].Nice != 0 {
		t.Errorf("parseProcStat: nice should be 0, but '%f'", *stat["cpu0"].Nice)
	}
	if *stat["cpu0"].System != 4570343 {
		t.Errorf("parseProcStat: system should be 4570343, but '%f'", *stat["cpu0"].System)
	}
	if *stat["cpu0"].Idle != 2196397805 {
		t.Errorf("parseProcStat: idle should be 2196397805, but '%f'", *stat["cpu0"].Idle)
	}
	if *stat["cpu0"].IoWait != 364717 {
		t.Errorf("parseProcStat: iowait should be 364717, but '%f'", *stat["cpu0"].IoWait)
	}
	if *stat["cpu0"].Irq != 4072 {
		t.Errorf("parseProcStat: irq should be 4072, but '%f'", *stat["cpu0"].Irq)
	}
	if *stat["cpu0"].SoftIrq != 1209169 {
		t.Errorf("parseProcStat: softirq should be 1209169, but '%f'", *stat["cpu0"].SoftIrq)
	}
	if *stat["cpu0"].Steal != 526216 {
		t.Errorf("parseProcStat: steal should be 526216, but '%f'", *stat["cpu0"].Steal)
	}
	if *stat["cpu0"].Guest != 1 {
		t.Errorf("parseProcStat: guest should be 1, but '%f'", *stat["cpu0"].Guest)
	}
	if *stat["cpu0"].GuestNice != 2 {
		t.Errorf("parseProcStat: guest should be 2, but '%f'", *stat["cpu0"].Guest)
	}
	if stat["cpu0"].Total != 2217859215 {
		t.Errorf("parseProcStat: total should be 2217859215, but '%f'", stat["cpu0"].Total)
	}
}

func TestParseProcStatsOldKernel(t *testing.T) {
	stab := `cpu0 14786890 0 4570343 2196397805`

	stat, _ := parseProcStat(stab)
	if len(stat) != 1 {
		t.Errorf("parseProcStat: size should be 1, but '%d'", len(stat))
	}
	if *stat["cpu0"].User != 14786890 {
		t.Errorf("parseProcStat: user should be 14786890, but '%f'", *stat["cpu0"].User)
	}
	if *stat["cpu0"].Nice != 0 {
		t.Errorf("parseProcStat: nice should be 0, but '%f'", *stat["cpu0"].Nice)
	}
	if *stat["cpu0"].System != 4570343 {
		t.Errorf("parseProcStat: system should be 4570343, but '%f'", *stat["cpu0"].System)
	}
	if *stat["cpu0"].Idle != 2196397805 {
		t.Errorf("parseProcStat: idle should be 2196397805, but '%f'", *stat["cpu0"].Idle)
	}
	if stat["cpu0"].IoWait != nil {
		t.Errorf("parseProcStat: iowait should be nil, but '%f'", *stat["cpu0"].IoWait)
	}
	if stat["cpu0"].Irq != nil {
		t.Errorf("parseProcStat: irq should be nil, but '%f'", *stat["cpu0"].Irq)
	}
	if stat["cpu0"].SoftIrq != nil {
		t.Errorf("parseProcStat: softirq should be nil, but '%f'", *stat["cpu0"].SoftIrq)
	}
	if stat["cpu0"].Steal != nil {
		t.Errorf("parseProcStat: steal should be nil, but '%f'", *stat["cpu0"].Steal)
	}
	if stat["cpu0"].Guest != nil {
		t.Errorf("parseProcStat: guest should be nil, but '%f'", *stat["cpu0"].Guest)
	}
}
