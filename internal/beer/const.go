package beer

import "github.com/ldassonville/beer-puller-api/pkg/api"

const EasyBeer = "EASY_BEER"
const LatencyBeer = "LATENCY_BEER"
const LazyBeer = "LAZY_BEER"
const FatalBeer = "FATAL_BEER"
const StolenBeer = "STOLEN_BEER"

var Catalog = []*api.Beer{

	{
		Code: EasyBeer,
		Name: "Easy beer",
		Desc: "Always fine beer. You will never been in trouble with it",
	},
	{
		Code: LatencyBeer,
		Name: "Latency beer",
		Desc: "Very popular beer, so it's take a long time to obtain it. You'll have to be patient to obtain your request.",
	},
	{
		Code: LazyBeer,
		Name: "Lazy beer",
		Desc: "This tasty beer takes time to chill. Therefore, it may not be ready immediately. ",
	},
	{
		Code: FatalBeer,
		Name: "Fatal beer",
		Desc: "With this beer it always ends badly. So it's up to you ",
	},
	{
		Code: StolenBeer,
		Name: "Stolen beer",
		Desc: "With this beer it always ends badly. So it's up to you ",
	},
}
