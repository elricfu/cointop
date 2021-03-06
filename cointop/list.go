package cointop

import (
	"time"

	types "github.com/miguelmota/cointop/pkg/api/types"
	"github.com/miguelmota/cointop/pkg/fcache"
)

func (ct *Cointop) updateCoins() error {
	list := []*coin{}
	cachekey := "allcoinsmap"

	var err error
	var allcoinsmap map[string]types.Coin
	cached, found := ct.cache.Get(cachekey)
	if found {
		// cache hit
		allcoinsmap, _ = cached.(map[string]types.Coin)
		ct.debuglog("soft cache hit")
	}

	// cache miss
	if allcoinsmap == nil {
		ct.debuglog("cache miss")
		allcoinsmap, err = ct.api.GetAllCoinData()
		if err != nil {
			return err
		}
		ct.cache.Set(cachekey, allcoinsmap, 10*time.Second)
		go func() {
			_ = fcache.Set(cachekey, allcoinsmap, 24*time.Hour)
		}()
	}

	if len(ct.allcoinsmap) == 0 {
		ct.allcoinsmap = map[string]*coin{}
	}
	for k, v := range allcoinsmap {
		last := ct.allcoinsmap[k]
		ct.allcoinsmap[k] = &coin{
			ID:               v.ID,
			Name:             v.Name,
			Symbol:           v.Symbol,
			Rank:             v.Rank,
			PriceUSD:         v.PriceUSD,
			PriceBTC:         v.PriceBTC,
			USD24HVolume:     v.USD24HVolume,
			MarketCapUSD:     v.MarketCapUSD,
			AvailableSupply:  v.AvailableSupply,
			TotalSupply:      v.TotalSupply,
			PercentChange1H:  v.PercentChange1H,
			PercentChange24H: v.PercentChange24H,
			PercentChange7D:  v.PercentChange7D,
			LastUpdated:      v.LastUpdated,
		}
		if last != nil {
			ct.allcoinsmap[k].Favorite = last.Favorite
		}
	}
	if len(ct.allcoins) == 0 {
		for i := range ct.allcoinsmap {
			coin := ct.allcoinsmap[i]
			list = append(list, coin)
		}
		ct.allcoins = list
		ct.sort(ct.sortby, ct.sortdesc, ct.allcoins)
	} else {
		// update list in place without changing order
		for i := range ct.allcoinsmap {
			cm := ct.allcoinsmap[i]
			for k := range ct.allcoins {
				c := ct.allcoins[k]
				if c.ID == cm.ID {
					// TODO: improve this
					c.ID = cm.ID
					c.Name = cm.Name
					c.Symbol = cm.Symbol
					c.Rank = cm.Rank
					c.PriceUSD = cm.PriceUSD
					c.PriceBTC = cm.PriceBTC
					c.USD24HVolume = cm.USD24HVolume
					c.MarketCapUSD = cm.MarketCapUSD
					c.AvailableSupply = cm.AvailableSupply
					c.TotalSupply = cm.TotalSupply
					c.PercentChange1H = cm.PercentChange1H
					c.PercentChange24H = cm.PercentChange24H
					c.PercentChange7D = cm.PercentChange7D
					c.LastUpdated = cm.LastUpdated
					c.Favorite = cm.Favorite
				}
			}
		}
	}
	return nil
}
