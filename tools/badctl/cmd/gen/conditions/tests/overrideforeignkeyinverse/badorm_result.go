// Code generated by badctl v0.0.0, DO NOT EDIT.
package overrideforeignkeyinverse

import badorm "github.com/ditrit/badaas/badorm"

func (m User) GetCreditCard() (*CreditCard, error) {
	return badorm.VerifyStructLoaded[CreditCard](&m.CreditCard)
}
