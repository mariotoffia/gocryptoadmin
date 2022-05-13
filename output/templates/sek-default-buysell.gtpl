-------------------------------------------------{{ translated . "separator" "Bought Price" "Fee" "SEK"}}{{ translated . "separator" "Sold Price" "Fee" "SEK"}}{{tax . "separator" "Tax" 30 "SEK"}}
|Bought    |Sold       |Pair    |Size           |{{ translated . "header" "Bought Price" "Fee" "SEK"}}{{ translated . "header" "Sold Price" "Fee" "SEK"}}{{tax . "header" "Tax" 30 "SEK"}}
-------------------------------------------------{{ translated . "separator" "Bought Price" "Fee" "SEK"}}{{ translated . "separator" "Sold Price" "Fee" "SEK"}}{{tax . "separator" "Tax" 30 "SEK"}}
{{range . }} 
{{- .GetBuy.GetCreatedAt.Format "2006-01-02 |" }}
{{- .GetSell.GetCreatedAt.Format "2006-01-02 " }}
{{- printf "|%-8s" .GetAssetPair.String }}
{{- printf "|% -15.8f|" .GetAssetSize }}
{{- translated .GetBuy "total-and-fee-positive" "Bought Price" "Fee" "SEK"}}
{{- translated .GetSell "total-and-fee-positive" "Sold Price" "Fee" "SEK"}}
{{- tax . "tax-all" "Tax" 30 "SEK"}}
{{end}}