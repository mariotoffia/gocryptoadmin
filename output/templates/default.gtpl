------------------------------------------------------------------------------------------------------{{ account . "separator" "optional"}}{{ translated . "separator"}}
|Exchange  |Side     |Tx Date             |Pair    |Price / Unit     |Fee          |Total Price      |{{ account . "header" "optional"}}{{ translated . "header"}}
------------------------------------------------------------------------------------------------------{{ account . "separator" "optional"}}{{ translated . "separator"}}
{{range . }} 
{{- printf "|%-10s" .GetExchange }}
{{- printf "|%-9v|" .GetSide }}
{{- .GetCreatedAt.Format "2006-01-02 15:04:05 " }}
{{- printf "|%-8s" .GetAssetPair.String }}
{{- printf "|%-17f" .GetPricePerUnit }}
{{- printf "|%-13f" .GetFee }}
{{- printf "|%-17f|" .GetTotalPrice }}
{{- account . "value" "optional"}}
{{- translated . "total-and-fee"}}
{{end}}