-------------------------------------------------------------------------------------------------------------------------------------------------------{{ account . "separator" "optional"}}{{ translated . "separator"}}
|Exchange  |Bought Date        |Sold Date           |Pair    |Size         |Bought           |Sold             |Bought Fee   |Sold Fee  |Tax          |{{ account . "header" "optional"}}{{ translated . "header"}}
-------------------------------------------------------------------------------------------------------------------------------------------------------{{ account . "separator" "optional"}}{{ translated . "separator"}}
{{range . }} 
{{- printf "|%-10s" .GetExchange }}
{{- .GetBuy.GetCreatedAt.Format "2006-01-02 15:04:05 |" }}
{{- .GetCreatedAt.Format "2006-01-02 15:04:05 " }}
{{- printf "|%-8s" .GetAssetPair.String }}
{{- printf "|% -13f" .GetAssetSize }}
{{- printf "|% -17f" .GetBuy.GetTotalPrice }}
{{- printf "|% -17f" .GetSell.GetTotalPrice }}
{{- printf "|% -13f|" .GetBuy.GetFee }}
{{- printf "|% -13f|" .GetSell.GetFee }}
{{- account . "value" "optional"}}
{{- translated . "total-and-fee"}}
{{end}}