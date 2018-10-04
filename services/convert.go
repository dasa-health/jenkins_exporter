package services

import "time"

// MillisToSecond convete milisegundos para segundos
func MillisToSecond(millis float64) float64 {
	return millis / 1000
}

// NanoTimestampToString converte um nano timestamp para uma data no padrao RFC3339
func NanoTimestampToString(value float64) string {
	utc, _ := time.LoadLocation("America/Sao_Paulo")
	tm := time.Unix((1538505491470 / 1000), 0).In(utc)

	return tm.Format(time.RFC3339)
}
