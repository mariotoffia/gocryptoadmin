---------------------------------------------------------------------------------------------------------------------------------------------{{tax . "separator" "Tax" 30 "EUR"}}
|Exchange  |Bought Date         |Sold Date           |Pair    |Size         |Purchase Price   |Selling Price    |Purchase Fee |Selling Fee  |{{tax . "header" "Tax" 30 "EUR"}}
---------------------------------------------------------------------------------------------------------------------------------------------{{tax . "separator" "Tax" 30 "EUR"}}
{{range . }} 
{{- printf "|%-10s|" .GetExchange }}
{{- .GetBuy.GetCreatedAt.Format "2006-01-02 15:04:05 |" }}
{{- .GetCreatedAt.Format "2006-01-02 15:04:05 " }}
{{- printf "|%-8s" .GetAssetPair.String }}
{{- printf "|% -13f" .GetAssetSize }}
{{- printf "|% -17f" .GetBuy.GetTotalPrice }}
{{- printf "|% -17f" .GetSell.GetTotalPrice }}
{{- printf "|% -13f" .GetBuy.GetFee }}
{{- printf "|% -13f|" .GetSell.GetFee }}
{{- tax . "tax-all" "Tax" 30 "EUR"}}
{{end}}