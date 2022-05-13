----------------------------------------------------------------------------------------------------------------------
|Exchange  |Side     |Tx Date             |Pair    |Size           |Price / Unit     |Fee          |Total Price      |
----------------------------------------------------------------------------------------------------------------------
{{range . }} 
{{- printf "|%-10s" .GetExchange }}
{{- printf "|%-9v|" .GetSide }}
{{- .GetCreatedAt.Format "2006-01-02 15:04:05 " }}
{{- printf "|%-8s" .GetAssetPair.String }}
{{- printf "|% -15.8f" .GetAssetSize }}
{{- printf "|% -17.2f" .GetPricePerUnit }} 
{{- printf "|% -13.2f" .GetFee }}
{{- printf "|% -17.2f|" .GetTotalPrice }}
{{end}}