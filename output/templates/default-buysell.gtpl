-------------------------------------------------------------------------------{{ translated . "separator" "Bought Price" "Fee" "EUR"}}{{ translated . "separator" "Sold Price" "Fee" "EUR"}}{{tax . "separator" "Tax" 30 "EUR"}}
|Exchange  |Bought Date         |Sold Date           |Pair    |Size           |{{ translated . "header" "Bought Price" "Fee" "EUR"}}{{ translated . "header" "Sold Price" "Fee" "EUR"}}{{tax . "header" "Tax" 30 "EUR"}}
-------------------------------------------------------------------------------{{ translated . "separator" "Bought Price" "Fee" "EUR"}}{{ translated . "separator" "Sold Price" "Fee" "EUR"}}{{tax . "separator" "Tax" 30 "EUR"}}
{{range . }} 
{{- printf "|%-10s|" .GetExchange }}
{{- .GetBuy.GetCreatedAt.Format "2006-01-02 15:04:05 |" }}
{{- .GetCreatedAt.Format "2006-01-02 15:04:05 " }}
{{- printf "|%-8s" .GetAssetPair.String }}
{{- printf "|% -15.8f|" .GetAssetSize }}
{{- translated .GetBuy "total-and-fee-positive" "Bought Price" "Fee" "EUR"}}
{{- translated .GetSell "total-and-fee-positive" "Sold Price" "Fee" "EUR"}}
{{- tax . "tax-all" "Tax" 30 "EUR"}}
{{end}}