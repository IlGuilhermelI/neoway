package dto

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/klassmann/cpfcnpj"
)

func trimSpace(value string) string {
	return strings.TrimRight(value, " ")
}
func formatToStringAndDateForInsert(value string) string {
	if strings.Contains(value, "NULL") {
		return "NULL"
	} else {
		return fmt.Sprintf("'%s'", value)
	}
}
func toDecimalFormat(value string) string {
	return strings.ReplaceAll(value, ",", ".")
}
func ValidateAllCpfAndCnpjLine(currentLine string) (result bool) {
	result = ValidateCpfAndCnpj(GetCpf(currentLine)) &&
		ValidateCpfAndCnpj(GetMostFrequentStore(currentLine)) &&
		ValidateCpfAndCnpj(GetLastPurchaseStore(currentLine))
	return result
}
func ValidateCpfAndCnpj(cpfOrCnpj string) bool {
	cpfOrCnpj = strings.ReplaceAll(cpfOrCnpj, "'", "")
	if cpfOrCnpj == "NULL" {
		return true
	}
	result := false
	if len(cpfOrCnpj) >= 14 {
		result = cpfcnpj.ValidateCNPJ(cpfOrCnpj)
	} else {
		result = cpfcnpj.ValidateCPF(cpfOrCnpj)
	}
	if !(result) {
		log.Print("The value: " + cpfOrCnpj + " is not a valid CPF or CNPJ, its line will not be included in the database")
	}
	return result

}
func InsertClientsPurchaseInformations(currentLine string) (query string) {
	query = fmt.Sprintf(`INSERT INTO CLIENTS_PURCHASE_INFORMATIONS
		(ID
		,CPF
		,PRIVATE
		,INCOMPLETE
		,LAST_PURCHASE_DATE
		,AVERAGE_TICKET
		,LAST_PURCHASE_TICKET
		,MOST_FREQUENT_STORE
		,LAST_PURCHASE_STORE)
		Values
		(DEFAULT
		,%s
		,%v
		,%v
		,%s
		,%s
		,%s
		,%s
		,%s`, GetCpf(currentLine), GetPrivate(currentLine), GetIncomplete(currentLine),
		GetLastPurchaseDate(currentLine), GetAverageTicket(currentLine), GetLastPurchaseTicket(currentLine),
		GetMostFrequentStore(currentLine), GetLastPurchaseStore(currentLine)+")")

	return query
}
func removeNonNumbersCaracters(cpfOrCnpj string) (result string) {
	result = strings.ReplaceAll(cpfOrCnpj, ".", "")
	result = strings.ReplaceAll(result, "/", "")
	return strings.ReplaceAll(result, "-", "")
}

func GetCpf(oneLine string) string {
	return formatToStringAndDateForInsert(removeNonNumbersCaracters(trimSpace(oneLine[0:14])))
}

func GetPrivate(oneLine string) bool {
	result, _ := strconv.ParseBool(oneLine[19:20])
	return result
}

func GetIncomplete(oneLine string) bool {
	result, _ := strconv.ParseBool(oneLine[31:32])
	return result
}

func GetLastPurchaseDate(oneLine string) string {
	return formatToStringAndDateForInsert(oneLine[43:53])
}

func GetAverageTicket(oneLine string) string {
	return toDecimalFormat(oneLine[65:87])
}

func GetLastPurchaseTicket(oneLine string) string {
	return toDecimalFormat(oneLine[87:110])
}

func GetMostFrequentStore(oneLine string) string {
	return formatToStringAndDateForInsert(removeNonNumbersCaracters(trimSpace(oneLine[111:129])))
}
func GetLastPurchaseStore(oneLine string) string {
	return formatToStringAndDateForInsert(removeNonNumbersCaracters(trimSpace(oneLine[131:])))
}
