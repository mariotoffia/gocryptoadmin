--------------------------------------------------------------------------------------------------------------------{{ account . "separator" "optional"}}{{ translated . "separator" "Total Price" "Fee"}}
|Exchange  |Side     |Tx Date             |Pair    |Size         |Price / Unit     |Fee          |Total Price      |{{ account . "header" "optional"}}{{ translated . "header" "Total Price" "Fee"}}
--------------------------------------------------------------------------------------------------------------------{{ account . "separator" "optional"}}{{ translated . "separator" "Total Price" "Fee"}}
{{range . }} 
{{- printf "|%-10s" .GetExchange }}
{{- printf "|%-9v|" .GetSide }}
{{- .GetCreatedAt.Format "2006-01-02 15:04:05 " }}
{{- printf "|%-8s" .GetAssetPair.String }}
{{- printf "|% -13f" .GetAssetSize }}
{{- printf "|% -17f" .GetPricePerUnit }}
{{- printf "|% -13f" .GetFee }}
{{- printf "|% -17f|" .GetTotalPrice }}
{{- account . "value" "optional"}}
{{- translated . "total-and-fee" "Total Price" "Fee"}}
{{end}}