package edsy

import (
	"fmt"
	"strings"

	"edDiscord/eddb"
)

type Stats struct {
	Cost                     int         `json:"cost"`
	CostVehicle              int         `json:"cost_vehicle"`
	CostRestock              int         `json:"cost_restock"`
	CostRearm                int         `json:"cost_rearm"`
	Mass                     float64     `json:"mass"`
	Pwrcap                   float64     `json:"pwrcap"`
	Pwrbst                   float64     `json:"pwrbst"`
	PwrdrawDep               []float64   `json:"pwrdraw_dep"`
	PwrdrawRet               []float64   `json:"pwrdraw_ret"`
	ThmloadCt                float64     `json:"thmload_ct"`
	ThmloadCfsd              float64     `json:"thmload_cfsd"`
	ThmloadHardpointWepfull  float64     `json:"thmload_hardpoint_wepfull"`
	ThmloadHardpointWepempty float64     `json:"thmload_hardpoint_wepempty"`
	ThmloadIscb              float64     `json:"thmload_iscb"`
	SpinupIscb               float64     `json:"spinup_iscb"`
	Jumpbst                  float64     `json:"jumpbst"`
	Fuelcap                  float64     `json:"fuelcap"`
	Cargocap                 float64     `json:"cargocap"`
	Cabincap                 float64     `json:"cabincap"`
	Scooprate                float64     `json:"scooprate"`
	Hullbst                  float64     `json:"hullbst"`
	Hullrnf                  float64     `json:"hullrnf"`
	Shieldbst                float64     `json:"shieldbst"`
	Shieldrnf                float64     `json:"shieldrnf"`
	Shieldrnfps              float64     `json:"shieldrnfps"`
	ShieldrnfpsAmmomax       float64     `json:"shieldrnfps_ammomax"`
	IntegImrp                float64     `json:"integ_imrp"`
	Dmgprot                  float64     `json:"dmgprot"`
	Dps                      float64     `json:"dps"`
	DpsAbs                   float64     `json:"dps_abs"`
	DpsThm                   float64     `json:"dps_thm"`
	DpsKin                   float64     `json:"dps_kin"`
	DpsExp                   float64     `json:"dps_exp"`
	DpsAxe                   float64     `json:"dps_axe"`
	DpsCau                   float64     `json:"dps_cau"`
	DpsNodistdraw            float64     `json:"dps_nodistdraw"`
	DpsDistdraw              float64     `json:"dps_distdraw"`
	AmmotimeWepcap           float64     `json:"ammotime_wepcap"`
	AmmotimeNocap            interface{} `json:"ammotime_nocap"`
	WepcapBurstCur           float64     `json:"wepcap_burst_cur"`
	WepcapBurstMax           float64     `json:"wepcap_burst_max"`
	WepchgSustainCur         float64     `json:"wepchg_sustain_cur"`
	WepchgSustainMax         float64     `json:"wepchg_sustain_max"`
	MassHull                 float64     `json:"mass_hull"`
	MassUnladen              float64     `json:"mass_unladen"`
	MassLaden                float64     `json:"mass_laden"`
	JumpLaden                float64     `json:"_jump_laden"`
	JumpUnladen              float64     `json:"_jump_unladen"`
	JumpMax                  float64     `json:"_jump_max"`
	RangeLaden               float64     `json:"_range_laden"`
	RangeUnladen             float64     `json:"_range_unladen"`
	Speed                    float64     `json:"_speed"`
	Boost                    float64     `json:"_boost"`
	Shields                  float64     `json:"_shields"`
	Skinres                  float64     `json:"_skinres"`
	Sthmres                  float64     `json:"_sthmres"`
	Sexpres                  float64     `json:"_sexpres"`
	Scaures                  float64     `json:"_scaures"`
	Armour                   float64     `json:"_armour"`
	Akinres                  float64     `json:"_akinres"`
	Athmres                  float64     `json:"_athmres"`
	Aexpres                  float64     `json:"_aexpres"`
	Acaures                  float64     `json:"_acaures"`
}

type PowerDist struct {
	Sys int `json:"sys"`
	Eng int `json:"eng"`
	Wep int `json:"wep"`
}

type ShipHull struct {
	Build     interface{} `json:"build"`
	Slotgroup string      `json:"slotgroup"`
	Slotnum   string      `json:"slotnum"`
	Hash      interface{} `json:"hash"`
	Modid     int         `json:"modid"`
	Module    struct {
		Fdid        int     `json:"fdid"`
		Fdname      string  `json:"fdname"`
		Eddbid      int     `json:"eddbid"`
		Id          int     `json:"id"`
		Stype       string  `json:"stype"`
		Name        string  `json:"name"`
		Class       int     `json:"class"`
		Cost        int     `json:"cost"`
		Retail      int     `json:"retail"`
		Topspd      int     `json:"topspd"`
		Bstspd      int     `json:"bstspd"`
		Mnv         int     `json:"mnv"`
		Shields     int     `json:"shields"`
		Armour      int     `json:"armour"`
		Mass        int     `json:"mass"`
		Fwdacc      float64 `json:"fwdacc"`
		Revacc      float64 `json:"revacc"`
		Latacc      float64 `json:"latacc"`
		Minthrust   float64 `json:"minthrust"`
		Boostcost   int     `json:"boostcost"`
		Pitch       int     `json:"pitch"`
		Yaw         int     `json:"yaw"`
		Roll        int     `json:"roll"`
		Pitchacc    int     `json:"pitchacc"`
		Yawacc      int     `json:"yawacc"`
		Rollacc     int     `json:"rollacc"`
		Minpitch    int     `json:"minpitch"`
		Heatcap     int     `json:"heatcap"`
		Heatdismin  float64 `json:"heatdismin"`
		Heatdismax  float64 `json:"heatdismax"`
		Fuelcost    int     `json:"fuelcost"`
		Fuelreserve float64 `json:"fuelreserve"`
		Hardness    int     `json:"hardness"`
		Masslock    int     `json:"masslock"`
		Crew        int     `json:"crew"`
		Slots       struct {
			Hardpoint []int `json:"hardpoint"`
			Utility   []int `json:"utility"`
			Component []int `json:"component"`
			Military  []int `json:"military"`
			Internal  []int `json:"internal"`
		} `json:"slots"`
		Slotnames struct {
			Internal []string `json:"internal"`
		} `json:"slotnames"`
		Stock struct {
			Hardpoint []int `json:"hardpoint"`
			Utility   []int `json:"utility"`
			Component []int `json:"component"`
			Military  []int `json:"military"`
			Internal  []int `json:"internal"`
		} `json:"stock"`
		Module struct {
			Field1 struct {
				Cost   int    `json:"cost"`
				Mass   int    `json:"mass"`
				Fdid   int    `json:"fdid"`
				Fdname string `json:"fdname"`
				Eddbid int    `json:"eddbid"`
			} `json:"40113"`
			Field2 struct {
				Cost   int    `json:"cost"`
				Mass   int    `json:"mass"`
				Fdid   int    `json:"fdid"`
				Fdname string `json:"fdname"`
				Eddbid int    `json:"eddbid"`
			} `json:"40114"`
			Field3 struct {
				Cost   int    `json:"cost"`
				Mass   int    `json:"mass"`
				Fdid   int    `json:"fdid"`
				Fdname string `json:"fdname"`
				Eddbid int    `json:"eddbid"`
			} `json:"40115"`
			Field4 struct {
				Cost   int    `json:"cost"`
				Mass   int    `json:"mass"`
				Fdid   int    `json:"fdid"`
				Fdname string `json:"fdname"`
				Eddbid int    `json:"eddbid"`
			} `json:"40122"`
			Field5 struct {
				Cost   int    `json:"cost"`
				Mass   int    `json:"mass"`
				Fdid   int    `json:"fdid"`
				Fdname string `json:"fdname"`
				Eddbid int    `json:"eddbid"`
			} `json:"40131"`
		} `json:"module"`
	} `json:"module"`
	Discounts    interface{} `json:"discounts"`
	Cost         int         `json:"cost"`
	Powered      bool        `json:"powered"`
	Priority     int         `json:"priority"`
	Bpid         int         `json:"bpid"`
	Bpgrade      int         `json:"bpgrade"`
	Bproll       int         `json:"bproll"`
	Expid        int         `json:"expid"`
	AttrModifier interface{} `json:"attrModifier"`
	AttrOverride interface{} `json:"attrOverride"`
}

type CargoHatch struct {
	Build     interface{} `json:"build"`
	Slotgroup string      `json:"slotgroup"`
	Slotnum   string      `json:"slotnum"`
	Hash      string      `json:"hash"`
	Modid     int         `json:"modid"`
	Module    struct {
		Mtype   string      `json:"mtype"`
		Cost    int         `json:"cost"`
		Name    string      `json:"name"`
		Class   int         `json:"class"`
		Rating  string      `json:"rating"`
		Pwrdraw float64     `json:"pwrdraw"`
		Fdid    interface{} `json:"fdid"`
		Fdname  string      `json:"fdname"`
		Eddbid  interface{} `json:"eddbid"`
	} `json:"module"`
	Discounts    interface{} `json:"discounts"`
	Cost         int         `json:"cost"`
	Powered      bool        `json:"powered"`
	Priority     int         `json:"priority"`
	Bpid         int         `json:"bpid"`
	Bpgrade      int         `json:"bpgrade"`
	Bproll       int         `json:"bproll"`
	Expid        int         `json:"expid"`
	AttrModifier interface{} `json:"attrModifier"`
	AttrOverride interface{} `json:"attrOverride"`
	Storedhash   string      `json:"storedhash"`
}

type Hardpoint struct {
	Mtype      string   `json:"mtype"`
	Cost       int      `json:"cost"`
	Name       string   `json:"name"`
	Mount      string   `json:"mount"`
	Class      int      `json:"class"`
	Rating     string   `json:"rating"`
	Mass       int      `json:"mass"`
	Integ      int      `json:"integ"`
	Pwrdraw    float64  `json:"pwrdraw"`
	Boottime   int      `json:"boottime"`
	Dps        float64  `json:"dps"`
	Damage     float64  `json:"damage"`
	Distdraw   float64  `json:"distdraw"`
	Thmload    float64  `json:"thmload"`
	Pierce     int      `json:"pierce"`
	Shotspd    int      `json:"shotspd,omitempty"`
	Rof        float64  `json:"rof"`
	Bstint     float64  `json:"bstint"`
	Ammoclip   int      `json:"ammoclip,omitempty"`
	Ammomax    int      `json:"ammomax,omitempty"`
	Rldtime    int      `json:"rldtime,omitempty"`
	Brcdmg     *float64 `json:"brcdmg"`
	Minbrc     int      `json:"minbrc"`
	Maxbrc     int      `json:"maxbrc"`
	Expwgt     int      `json:"expwgt,omitempty"`
	Ammocost   int      `json:"ammocost,omitempty"`
	Fdid       int      `json:"fdid"`
	Fdname     string   `json:"fdname"`
	Eddbid     int      `json:"eddbid"`
	Tag        string   `json:"tag,omitempty"`
	Duration   float64  `json:"duration,omitempty"`
	Dmgmul     int      `json:"dmgmul,omitempty"`
	Maximumrng int      `json:"maximumrng,omitempty"`
	Abswgt     int      `json:"abswgt,omitempty"`
	Axewgt     int      `json:"axewgt,omitempty"`
	Dmgfall    int      `json:"dmgfall,omitempty"`
	Limit      string   `json:"limit,omitempty"`
	Thmwgt     int      `json:"thmwgt,omitempty"`
}

type HardpointSlot struct {
	Build        interface{} `json:"build"`
	Slotgroup    string      `json:"slotgroup"`
	Slotnum      int         `json:"slotnum"`
	Hash         *string     `json:"hash"`
	Modid        int         `json:"modid"`
	Module       *Hardpoint  `json:"module"`
	Discounts    interface{} `json:"discounts"`
	Cost         int         `json:"cost"`
	Powered      bool        `json:"powered"`
	Priority     int         `json:"priority"`
	Bpid         interface{} `json:"bpid"`
	Bpgrade      int         `json:"bpgrade"`
	Bproll       int         `json:"bproll"`
	Expid        interface{} `json:"expid"`
	AttrModifier *struct {
		Pwrdraw  float64 `json:"pwrdraw"`
		Damage   float64 `json:"damage"`
		Distdraw float64 `json:"distdraw"`
		Thmload  float64 `json:"thmload"`
	} `json:"attrModifier"`
	AttrOverride *struct {
	} `json:"attrOverride"`
	Storedhash string `json:"storedhash,omitempty"`
}

func (h HardpointSlot) String() string {
	if h.Module == nil {
		return ""
	}
	mount := ""
	switch h.Module.Mount {
	case "T":
		mount = "Turret"
	case "G":
		mount = "Gimballed"
	case "F":
		mount = "Fixed"

	}
	return fmt.Sprintf(
		"[**%v%v %v %s**](%v)",
		h.Module.Class,
		h.Module.Rating,
		mount,
		h.Module.Name,
		eddb.GetEDDBModule(h.Module.Eddbid),
	)
}

type Utility struct {
	Mtype      string  `json:"mtype"`
	Cost       int     `json:"cost"`
	Name       string  `json:"name"`
	Class      int     `json:"class"`
	Rating     string  `json:"rating"`
	Mass       float64 `json:"mass"`
	Integ      int     `json:"integ"`
	Pwrdraw    float64 `json:"pwrdraw"`
	Passive    int     `json:"passive,omitempty"`
	Boottime   int     `json:"boottime"`
	Distdraw   int     `json:"distdraw,omitempty"`
	Thmload    float64 `json:"thmload,omitempty"`
	Rof        int     `json:"rof,omitempty"`
	Bstint     float64 `json:"bstint,omitempty"`
	Ammoclip   int     `json:"ammoclip,omitempty"`
	Ammomax    int     `json:"ammomax,omitempty"`
	Rldtime    float64 `json:"rldtime,omitempty"`
	Jamdur     int     `json:"jamdur,omitempty"`
	Ammocost   int     `json:"ammocost,omitempty"`
	Fdid       int     `json:"fdid"`
	Fdname     string  `json:"fdname"`
	Eddbid     int     `json:"eddbid"`
	Mount      string  `json:"mount,omitempty"`
	Dps        int     `json:"dps,omitempty"`
	Damage     float64 `json:"damage,omitempty"`
	Maximumrng int     `json:"maximumrng,omitempty"`
	Shotspd    int     `json:"shotspd,omitempty"`
	Bstrof     int     `json:"bstrof,omitempty"`
	Bstsize    int     `json:"bstsize,omitempty"`
	Jitter     float64 `json:"jitter,omitempty"`
	Kinwgt     int     `json:"kinwgt,omitempty"`
	Scanrng    int     `json:"scanrng,omitempty"`
	Maxangle   int     `json:"maxangle,omitempty"`
	Scantime   int     `json:"scantime,omitempty"`
	Limit      string  `json:"limit,omitempty"`
}

type UtilitySlot struct {
	Build        interface{} `json:"build"`
	Slotgroup    string      `json:"slotgroup"`
	Slotnum      int         `json:"slotnum"`
	Hash         *string     `json:"hash"`
	Modid        int         `json:"modid"`
	Module       *Utility    `json:"module"`
	Discounts    interface{} `json:"discounts"`
	Cost         int         `json:"cost"`
	Powered      bool        `json:"powered"`
	Priority     int         `json:"priority"`
	Bpid         interface{} `json:"bpid"`
	Bpgrade      int         `json:"bpgrade"`
	Bproll       int         `json:"bproll"`
	Expid        int         `json:"expid"`
	AttrModifier *struct {
		Mass    float64 `json:"mass"`
		Ammomax float64 `json:"ammomax,omitempty"`
		Rldtime float64 `json:"rldtime,omitempty"`
		Integ   int     `json:"integ,omitempty"`
	} `json:"attrModifier"`
	AttrOverride *struct {
	} `json:"attrOverride"`
	Storedhash string `json:"storedhash,omitempty"`
}

func (u UtilitySlot) String() string {
	if u.Module != nil {
		return ""
	}

	return fmt.Sprintf(
		"[**%v%v %s**](%v)",
		u.Module.Class,
		u.Module.Rating,
		u.Module.Name,
		eddb.GetEDDBModule(u.Module.Eddbid),
	)
}

type CoreModule struct {
	Mtype      string  `json:"mtype"`
	Cost       int     `json:"cost"`
	Name       string  `json:"name"`
	Class      int     `json:"class"`
	Rating     string  `json:"rating"`
	Mass       int     `json:"mass,omitempty"`
	Hullbst    int     `json:"hullbst,omitempty"`
	Kinres     float64 `json:"kinres,omitempty"`
	Thmres     float64 `json:"thmres,omitempty"`
	Expres     float64 `json:"expres,omitempty"`
	Axeres     float64 `json:"axeres,omitempty"`
	Fdid       int     `json:"fdid"`
	Fdname     string  `json:"fdname"`
	Eddbid     int     `json:"eddbid"`
	Integ      int     `json:"integ,omitempty"`
	Pwrcap     float64 `json:"pwrcap,omitempty"`
	Heateff    float64 `json:"heateff,omitempty"`
	Pwrdraw    float64 `json:"pwrdraw,omitempty"`
	Boottime   int     `json:"boottime,omitempty"`
	Engminmass int     `json:"engminmass,omitempty"`
	Engoptmass int     `json:"engoptmass,omitempty"`
	Engmaxmass int     `json:"engmaxmass,omitempty"`
	Engminmul  int     `json:"engminmul,omitempty"`
	Engoptmul  int     `json:"engoptmul,omitempty"`
	Engmaxmul  int     `json:"engmaxmul,omitempty"`
	Engheat    float64 `json:"engheat,omitempty"`
	Fsdoptmass int     `json:"fsdoptmass,omitempty"`
	Fsdheat    int     `json:"fsdheat,omitempty"`
	Maxfuel    float64 `json:"maxfuel,omitempty"`
	Fuelmul    float64 `json:"fuelmul,omitempty"`
	Fuelpower  float64 `json:"fuelpower,omitempty"`
	Emgcylife  int     `json:"emgcylife,omitempty"`
	Wepcap     int     `json:"wepcap,omitempty"`
	Wepchg     float64 `json:"wepchg,omitempty"`
	Engcap     int     `json:"engcap,omitempty"`
	Engchg     float64 `json:"engchg,omitempty"`
	Syscap     int     `json:"syscap,omitempty"`
	Syschg     float64 `json:"syschg,omitempty"`
	Maxrng     int     `json:"maxrng,omitempty"`
	Scanangle  int     `json:"scanangle,omitempty"`
	Typemis    int     `json:"typemis,omitempty"`
	Fuelcap    int     `json:"fuelcap,omitempty"`
}

type CoreSlot struct {
	Build        interface{} `json:"build"`
	Slotgroup    string      `json:"slotgroup"`
	Slotnum      int         `json:"slotnum"`
	Hash         string      `json:"hash"`
	Modid        int         `json:"modid"`
	Module       CoreModule  `json:"module"`
	Discounts    interface{} `json:"discounts"`
	Cost         int         `json:"cost"`
	Powered      bool        `json:"powered"`
	Priority     int         `json:"priority"`
	Bpid         interface{} `json:"bpid"`
	Bpgrade      int         `json:"bpgrade"`
	Bproll       int         `json:"bproll"`
	Expid        interface{} `json:"expid"`
	AttrModifier *struct {
		Kinres     float64 `json:"kinres,omitempty"`
		Thmres     float64 `json:"thmres,omitempty"`
		Expres     float64 `json:"expres,omitempty"`
		Mass       float64 `json:"mass,omitempty"`
		Integ      float64 `json:"integ,omitempty"`
		Pwrcap     float64 `json:"pwrcap,omitempty"`
		Heateff    float64 `json:"heateff,omitempty"`
		Pwrdraw    float64 `json:"pwrdraw,omitempty"`
		Engoptmass float64 `json:"engoptmass,omitempty"`
		Engoptmul  float64 `json:"engoptmul,omitempty"`
		Engheat    float64 `json:"engheat,omitempty"`
		Boottime   float64 `json:"boottime,omitempty"`
		Fsdoptmass float64 `json:"fsdoptmass,omitempty"`
		Fsdheat    float64 `json:"fsdheat,omitempty"`
		Wepcap     float64 `json:"wepcap,omitempty"`
		Wepchg     float64 `json:"wepchg,omitempty"`
		Engcap     float64 `json:"engcap,omitempty"`
		Engchg     float64 `json:"engchg,omitempty"`
		Syscap     float64 `json:"syscap,omitempty"`
		Syschg     float64 `json:"syschg,omitempty"`
		Scanangle  float64 `json:"scanangle,omitempty"`
	} `json:"attrModifier"`
	AttrOverride *struct {
	} `json:"attrOverride"`
	Storedhash string `json:"storedhash"`
}

func (c CoreSlot) String() string {
	return fmt.Sprintf(
		"[%d%s %s](%v)",
		c.Module.Class,
		c.Module.Rating,
		c.Module.Name,
		eddb.GetEDDBModule(c.Module.Eddbid),
	)
}

type Optional struct {
	Mtype        string  `json:"mtype"`
	Cost         int     `json:"cost"`
	Name         string  `json:"name"`
	Class        int     `json:"class"`
	Rating       string  `json:"rating"`
	Cargocap     int     `json:"cargocap,omitempty"`
	Fdid         int     `json:"fdid"`
	Fdname       string  `json:"fdname"`
	Eddbid       int     `json:"eddbid"`
	Mass         int     `json:"mass,omitempty"`
	Integ        int     `json:"integ,omitempty"`
	Pwrdraw      float64 `json:"pwrdraw,omitempty"`
	Boottime     int     `json:"boottime,omitempty"`
	Genminmass   int     `json:"genminmass,omitempty"`
	Genoptmass   int     `json:"genoptmass,omitempty"`
	Genmaxmass   int     `json:"genmaxmass,omitempty"`
	Genminmul    int     `json:"genminmul,omitempty"`
	Genoptmul    int     `json:"genoptmul,omitempty"`
	Genmaxmul    int     `json:"genmaxmul,omitempty"`
	Genrate      float64 `json:"genrate,omitempty"`
	Bgenrate     float64 `json:"bgenrate,omitempty"`
	Genpwr       float64 `json:"genpwr,omitempty"`
	Kinres       float64 `json:"kinres,omitempty"`
	Thmres       float64 `json:"thmres,omitempty"`
	Expres       float64 `json:"expres,omitempty"`
	Axeres       float64 `json:"axeres,omitempty"`
	Limit        string  `json:"limit,omitempty"`
	Tag          string  `json:"tag,omitempty"`
	Shieldrnf    int     `json:"shieldrnf,omitempty"`
	Noblueprints struct {
		Field1 int `json:"*"`
	} `json:"noblueprints,omitempty"`
}

type OptionalSlot struct {
	Build        interface{} `json:"build"`
	Slotgroup    string      `json:"slotgroup"`
	Slotnum      int         `json:"slotnum"`
	Hash         *string     `json:"hash"`
	Modid        int         `json:"modid"`
	Module       *Optional   `json:"module"`
	Discounts    interface{} `json:"discounts"`
	Cost         int         `json:"cost"`
	Powered      bool        `json:"powered"`
	Priority     int         `json:"priority"`
	Bpid         interface{} `json:"bpid"`
	Bpgrade      int         `json:"bpgrade"`
	Bproll       int         `json:"bproll"`
	Expid        interface{} `json:"expid"`
	AttrModifier *struct {
		Mass       float64 `json:"mass"`
		Integ      float64 `json:"integ"`
		Pwrdraw    float64 `json:"pwrdraw"`
		Genoptmass float64 `json:"genoptmass"`
		Genoptmul  float64 `json:"genoptmul"`
	} `json:"attrModifier"`
	AttrOverride *struct {
	} `json:"attrOverride"`
	Storedhash string `json:"storedhash,omitempty"`
}

func (o OptionalSlot) String() string {
	if o.Module != nil {
		return ""
	}
	shortName := strings.TrimRight(strings.Split(o.Module.Name, "(")[0], " ")
	return fmt.Sprintf(
		"[**%v%v %s**](%v)",
		o.Module.Class,
		o.Module.Rating,
		shortName,
		eddb.GetEDDBModule(o.Module.Eddbid),
	)
}

type BuildSlots struct {
	Ship struct {
		Hull  ShipHull   `json:"hull"`
		Hatch CargoHatch `json:"hatch"`
	} `json:"ship"`
	Hardpoint []HardpointSlot `json:"hardpoint"`
	Utility   []UtilitySlot   `json:"utility"`
	Component CoreSlots       `json:"component"`
	Military  []OptionalSlot  `json:"military"`
	Internal  []OptionalSlot  `json:"internal"`
}

type CoreSlots []CoreSlot

func (i CoreSlots) String() string {
	return fmt.Sprintf(
		"**%v**\n**%v**\n**%v**\n**%v**\n**%v**\n**%v**\n**%v**\n**%v**\n",
		i[0],
		i[1],
		i[2],
		i[3],
		i[4],
		i[5],
		i[6],
		i[7],
	)
}

type Build struct {
	Hash      string      `json:"hash"`
	Stats     Stats       `json:"stats"`
	Shipid    int         `json:"shipid"`
	Name      string      `json:"name"`
	Nametag   string      `json:"nametag"`
	InaraAcct interface{} `json:"inaraAcct"`
	InaraShip interface{} `json:"inaraShip"`
	Crewdist  PowerDist   `json:"crewdist"`
	Powerdist PowerDist   `json:"powerdist"`
	Slots     BuildSlots  `json:"slots"`
}

func (b *Build) GetEDDBLink() string {
	shipID := b.Slots.Ship.Hull.Module.Eddbid
	moduleID := make(
		[]int, len(b.Slots.Hardpoint)+
			len(b.Slots.Component)+len(b.Slots.Military)+len(b.Slots.Utility)+len(b.Slots.Internal),
	)
	for _, slot := range b.Slots.Hardpoint {
		if slot.Module != nil {
			moduleID = append(moduleID, slot.Module.Eddbid)
		}

	}
	for _, slot := range b.Slots.Component {
		moduleID = append(moduleID, slot.Module.Eddbid)
	}
	for _, slot := range b.Slots.Military {
		if slot.Module != nil {
			moduleID = append(moduleID, slot.Module.Eddbid)
		}
	}
	for _, slot := range b.Slots.Utility {
		if slot.Module != nil {
			moduleID = append(moduleID, slot.Module.Eddbid)
		}
	}
	for _, slot := range b.Slots.Internal {
		if slot.Module != nil {
			moduleID = append(moduleID, slot.Module.Eddbid)
		}
	}

	return eddb.GetEDDBLink(shipID, moduleID)
}
