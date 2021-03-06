package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"strings"
    "crypto/md5"
    "encoding/hex"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var ptyPrefix = "pty:"
var accountPrefix = "acct:"
var accountsKey = "accounts"

type PTY struct {
	CUSIP		string 	   `json:"cusip"`
	Name		string 	   `json:"name"`
    AdrStreet   string     `json:"adrStreet"`
    AdrCity     string     `json:"adrCity"`
    AdrPostcode string     `json:"adrPostcode"`
    AdrState    string     `json:"adrState"`
    BuyValue    float64    `json:"buyval"`
    MktValue    float64    `json:"mktval"`
    Qty         int        `json:"quantity"`
    Owners      []Owner    `json:"owner"`
    PT4Sale     []ForSale  `json:"forsale"`
    Renters     []Renter   `json:"renters"`
    Links       []UrlLnk   `json:"urlLink"`
    Rent        float64    `json:"rent"`
    Issuer      string     `json:"issuer"`
    IssueDate   string     `json:"issueDate"`
    Status      string     `json:"status"`
}

type Owner struct {
	InvestorID string    `json:"invid"`
	Quantity int      `json:"quantity"`
}

type Renter struct {
    RenterID string    `json:"rentid"`
}

type ForSale struct {
    InvestorID string   `json:"invid"`
    Quantity   int      `json:"quantity"`
    SellVal    float64  `json:"sellval"`
}

type UrlLnk struct {
    Url         string   `json:"url"`
    UrlType     string   `json:"urlType"`
}

type Transaction struct {
	CUSIP       string   `json:"cusip"`
	FromCompany string   `json:"fromCompany"`
	ToCompany   string   `json:"toCompany"`
	Quantity    int      `json:"quantity"`
}

type AddForSale struct {
    CUSIP       string   `json:"cusip"`
    FromCompany string   `json:"fromCompany"`
    Quantity    int      `json:"quantity"`
    SellVal     float64  `json:"sellval"`
}

type Account struct {
	ID          string  `json:"id"`
	Prefix      string  `json:"prefix"`
    CashBalance float64 `json:"cashBalance"`
	AssetsIds   []string `json:"assetIds"`
    RentingPty  string   `json:"rentingpty"`
}

type SetRenter struct {
    CUSIP       string  `json:"cusip"`
    Action      string  `json:"action"`
    RenterName  string  `json:"invid"`
}

type SetRentValue struct {
    CUSIP       string  `json:"cusip"`
    Value       float64 `json:"value"`
    Issuer      string  `json:"invid"`
}

type UpdateMktVal struct {
    CUSIP       string   `json:"cusip"`
    MktValue    float64  `json:"mktval"`
}


type PayRent struct {
    CUSIP       string   `json:"cusip"`
    Payment     float64  `json:"payment"`
    Issuer      string   `json:"issuer"`
}

type SimpleChaincode struct {
}

const (
    millisPerSecond     = int64(time.Second / time.Millisecond)
    nanosPerMillisecond = int64(time.Millisecond / time.Nanosecond)
)

func msToTime(ms string) (time.Time, error) {
    msInt, err := strconv.ParseInt(ms, 10, 64)
    if err != nil {
        return time.Time{}, err
    }

    return time.Unix(msInt/millisPerSecond,
        (msInt%millisPerSecond)*nanosPerMillisecond), nil
}

func genHash(text string) (string, error) {

    hasher := md5.New()
    hasher.Write([]byte(strings.ToUpper(text)))


    // maturityDate := t.AddDate(0, 0, days)
    // month := int(maturityDate.Month())
    // day := maturityDate.Day()

    suffix := hex.EncodeToString(hasher.Sum(nil))
    return suffix, nil

}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    // Initialize the collection of commercial paper keys
    fmt.Println("Initializing Property keys collection")
	
    // Check if state already exists
    fmt.Println("Getting Property Keys")
    keysBytes, err := stub.GetState("PtyKeys")
    if keysBytes == nil {
        fmt.Println("Cannot find PtyKeys, will reinitialize everything")
        var blank []string
        blankBytes, _ := json.Marshal(&blank)
        err := stub.PutState("PtyKeys", blankBytes)
        if err != nil {
            fmt.Println("Failed to initialize property key collection")
        }
    } else if err != nil {
         fmt.Println("Failed to initialize property key collection")
    } else {
        fmt.Println("Found property keyBytes. Will not overwrite keys.")
    }



	fmt.Println("Initialization complete")

	return nil, nil
}

func (t *SimpleChaincode) createAccounts(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

    //                  0
    // "number of accounts to create"
    var err error
    numAccounts, err := strconv.Atoi(args[0])
    if err != nil {
        fmt.Println("error creating accounts with input")
        return nil, errors.New("createAccounts accepts a single integer argument")
    }
    //create a bunch of accounts
    var account Account
    counter := 1
    for counter <= numAccounts {
        var prefix string
        suffix := "000A"
        if counter < 10 {
            prefix = strconv.Itoa(counter) + "0" + suffix
        } else {
            prefix = strconv.Itoa(counter) + suffix
        }
        var assetIds []string
        account = Account{ID: "company" + strconv.Itoa(counter), Prefix: prefix, CashBalance: 10000000.0, AssetsIds: assetIds}
        accountBytes, err := json.Marshal(&account)
        if err != nil {
            fmt.Println("error creating account" + account.ID)
            return nil, errors.New("Error creating account " + account.ID)
        }
        err = stub.PutState(accountPrefix+account.ID, accountBytes)
        counter++
        fmt.Println("created account" + accountPrefix + account.ID)
    }

    fmt.Println("Accounts created")
    return nil, nil

}

func (t *SimpleChaincode) createAccount(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    // Obtain the username to associate with the account
    if len(args) != 1 {
        fmt.Println("Error obtaining username")
        return nil, errors.New("createAccount accepts a single username argument")
    }
    username := args[0]
    fmt.Println(username)
    fmt.Println("thats the username!")
    // Build an account object for the user
    var assetIds []string
    suffix := "000A"
    prefix := username + suffix
    var account = Account{ID: username, Prefix: prefix, CashBalance: 10000000.0, AssetsIds: assetIds}
    accountBytes, err := json.Marshal(&account)
    fmt.Println("Creating accounts")
    if err != nil {
        fmt.Println("error creating account" + account.ID)
        return nil, errors.New("Error creating account " + account.ID)
    }
    
    fmt.Println("Attempting to get state of any existing account for " + account.ID)
    existingBytes, err := stub.GetState(accountPrefix + account.ID)
	if err == nil {
        
        var company Account
        err = json.Unmarshal(existingBytes, &company)
        if err != nil {
            fmt.Println("Error unmarshalling account " + account.ID + "\n--->: " + err.Error())
            
            if strings.Contains(err.Error(), "unexpected end") {
                fmt.Println("No data means existing account found for " + account.ID + ", initializing account.")
                err = stub.PutState(accountPrefix+account.ID, accountBytes)
                
                if err == nil {
                    fmt.Println("created account" + accountPrefix + account.ID)
                    return nil, nil
                } else {
                    fmt.Println("failed to create initialize account for " + account.ID)
                    return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
                }
            } else {
                return nil, errors.New("Error unmarshalling existing account " + account.ID)
            }
        } else {
            fmt.Println("Account already exists for " + account.ID + " " + company.ID)
		    return nil, errors.New("Can't reinitialize existing user " + account.ID)
        }
    } else {
        
        fmt.Println("No existing account found for " + account.ID + ", initializing account.")
        err = stub.PutState(accountPrefix+account.ID, accountBytes)
        
        if err == nil {
            fmt.Println("created account" + accountPrefix + account.ID)
            return nil, nil
        } else {
            fmt.Println("failed to create initialize account for " + account.ID)
            return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
        }
        
    }
    
    
}

func (t *SimpleChaincode) updateMktVal(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    if len(args) != 1 {
        fmt.Println("error invalid arguments")
        return nil, errors.New("Incorrect number of arguments. Expecting commercial paper record")
    }

    /*
        type UpdateMktVal struct {
        CUSIP       string   `json:"cusip"`
        MktValue    float64  `json:"mktval"`
}   */

    var cp UpdateMktVal
    var err error

    var newstring = args[0]
    newstring = strings.Replace(args[0],"'","\"",-1)

    fmt.Println("Unmarshalling CP")
    err = json.Unmarshal([]byte(newstring), &cp)
    if err != nil {
        fmt.Println("error invalid paper issue")
        fmt.Println("error: ",err)
        return nil, errors.New("Invalid commercial paper issue")
    }


    fmt.Println("Getting State on CP " + cp.CUSIP)
    cpRxBytes, err := stub.GetState(ptyPrefix+cp.CUSIP)

    if cpRxBytes != nil {
        fmt.Println("CUSIP exists")
        
        var cprx PTY
        fmt.Println("Unmarshalling CP " + cp.CUSIP)
        err = json.Unmarshal(cpRxBytes, &cprx)
        if err != nil {
            fmt.Println("Error unmarshalling cp " + cp.CUSIP)
            return nil, errors.New("Error unmarshalling cp " + cp.CUSIP)
        }

        cprx.MktValue = cp.MktValue
        cprx.Status = "Approved"

        cpWriteBytes, err := json.Marshal(&cprx)
        if err != nil {
            fmt.Println("Error marshalling cp")
            return nil, errors.New("Error issuing commercial paper")
        }
        err = stub.PutState(ptyPrefix+cp.CUSIP, cpWriteBytes)
        if err != nil {
            fmt.Println("Error issuing paper")
            return nil, errors.New("Error issuing commercial paper")
        }

        fmt.Println("Updated commercial paper %+v\n", cprx)
        return nil, nil
    } else {
        return nil, errors.New("Could not find Property Token")
    }

}

func (t *SimpleChaincode) setRent(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    if len(args) != 1 {
        fmt.Println("error invalid arguments")
        return nil, errors.New("Incorrect number of arguments. Expecting commercial paper record")
    }

/*type SetRentValue struct {
    CUSIP       string  `json:"cusip"`
    Value       float64 `json:"value"`
    Issuer      string  `json:"invid"`
}*/

    var cp SetRentValue
    var err error

    var newstring = args[0]
    newstring = strings.Replace(args[0],"'","\"",-1)

    fmt.Println("Unmarshalling Data")
    err = json.Unmarshal([]byte(newstring), &cp)
    if err != nil {
        fmt.Println("error invalid Data issue")
        fmt.Println("error: ",err)
        return nil, errors.New("Invalid Data issue")
    }
    fmt.Println("Getting state of - " + accountPrefix + cp.Issuer)
    accountBytes, err := stub.GetState(accountPrefix + cp.Issuer)
    if err != nil {
        fmt.Println("Error Getting state of - " + accountPrefix + cp.Issuer)
        return nil, errors.New("Error retrieving account " + cp.Issuer)
    }
    if accountBytes == nil {
        fmt.Println("Lol how did you get here")
        return nil, errors.New("Error retrieving account " + cp.Issuer)
    }

    fmt.Println("Getting State on PTY " + cp.CUSIP)
    cpRxBytes, err := stub.GetState(ptyPrefix+cp.CUSIP)

    if cpRxBytes != nil {
        fmt.Println("CUSIP exists")
        
        var cprx PTY
        fmt.Println("Unmarshalling PTY " + cp.CUSIP)
        err = json.Unmarshal(cpRxBytes, &cprx)
        if err != nil {
            fmt.Println("Error unmarshalling cp " + cp.CUSIP)
            return nil, errors.New("Error unmarshalling cp " + cp.CUSIP)
        }

        cprx.Rent = cp.Value

        cpWriteBytes, err := json.Marshal(&cprx)
        if err != nil {
            fmt.Println("Error marshalling cp")
            return nil, errors.New("Error issuing commercial paper")
        }
        err = stub.PutState(ptyPrefix+cp.CUSIP, cpWriteBytes)
        if err != nil {
            fmt.Println("Error issuing paper")
            return nil, errors.New("Error issuing commercial paper")
        }

        fmt.Println("Updated commercial paper %+v\n", cprx)
        return nil, nil
    } else {
        return nil, errors.New("Could not find Property Token")
    }

}

func (t *SimpleChaincode) processRent(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    if len(args) != 1 {
        fmt.Println("error invalid arguments")
        return nil, errors.New("Incorrect number of arguments. Expecting payRent record")
    }

    /*
        type UpdateMktVal struct {
        CUSIP       string   `json:"cusip"`
        MktValue    float64  `json:"mktval"`
}   */

    var cp PayRent
    var err error

    var newstring = args[0]
    newstring = strings.Replace(args[0],"'","\"",-1)

    fmt.Println("Unmarshalling CP")
    err = json.Unmarshal([]byte(newstring), &cp)
    if err != nil {
        fmt.Println("error invalid paper issue")
        fmt.Println("error: ",err)
        return nil, errors.New("Invalid commercial paper issue")
    }
    var username = cp.Issuer
 //   suffix := "000A"
//    prefix := username + suffix
    var renter Account

    // Get state of renter account
    existingBytes, err := stub.GetState(accountPrefix + username)
    if err == nil {
        err = json.Unmarshal(existingBytes, &renter)
        if err == nil {

            // Check he has enough cash

            if (renter.CashBalance >= cp.Payment) {
                renter.CashBalance -= t.calcRent(stub, cp.CUSIP);
            } else {
                fmt.Println("Renter doesn't have enough money!")
                return nil, errors.New("Renter doens't have enough money!")
            }

        } else {
            fmt.Println("Cannot find renter account")
            return nil, errors.New("Failed to find renter account")
        }
    } else {
        fmt.Println("Unable to get account information")
        return nil, errors.New("Failed to get account information")
    }
    var currOwners []Owner
    var currSellers []ForSale
    var rentPerToken float64    
    var cprx PTY
    // Get state of the PTY that rent is being paid out to.

    fmt.Println("Getting State on PTY " + cp.CUSIP)
    cpRxBytes, err := stub.GetState(ptyPrefix+cp.CUSIP)

    if cpRxBytes != nil {
        fmt.Println("CUSIP exists")
        
        
        fmt.Println("Unmarshalling CP " + cp.CUSIP)
        err = json.Unmarshal(cpRxBytes, &cprx)
        if err != nil {
            fmt.Println("Error unmarshalling cp " + cp.CUSIP)
            return nil, errors.New("Error unmarshalling cp " + cp.CUSIP)
        }

        // Add in logic to figure out quantities each owner has, divide by quantity and send to all owners

        currOwners = cprx.Owners
        currSellers = cprx.PT4Sale
        // Calculate what each token gets in terms of rent

        rentPerToken = t.calcRent(stub, cp.CUSIP)/float64(cprx.Qty)

        // Making sure that we calculate what is up for sale as well as part of quantity
        
    }
    for _, owner := range currSellers {
            // for _, curOwner := range currOwners {
            //     if owner.InvestorID == curOwner.InvestorID {
            //         curOwner.Quantity += owner.Quantity
            //         fmt.Println("Found owner that has more quantity, adding to curOnwer")
            //         fmt.Println(curOwner.Quantity)
            //     }
            // }

        for i := 0; i < len(currOwners); i++ {
            if owner.InvestorID == currOwners[i].InvestorID {
                currOwners[i].Quantity += owner.Quantity
                fmt.Println("Found owner that has more quantity, adding to curOnwer")
                fmt.Println(currOwners[i].Quantity)
            }
        }
    }
    
    for _, curOwner := range currOwners {
        existingBytes, err := stub.GetState(accountPrefix + curOwner.InvestorID)
        if err == nil {
            // Unmarshal the damn bytes
            var ownerAccts Account
            err = json.Unmarshal(existingBytes,&ownerAccts)
            fmt.Println("curOwner quantity is: ", curOwner.Quantity)
            ownerAccts.CashBalance+=rentPerToken * float64(curOwner.Quantity)
            acctBytes, err := json.Marshal(&ownerAccts)
            if err == nil {
                err = stub.PutState(accountPrefix+curOwner.InvestorID, acctBytes)
            } else {
                fmt.Println("Egads, something went wrong with Marshalling")
                return nil, errors.New("Egads, something went wrong with Marshalling")
            }

            renterBytes, err := json.Marshal(&renter)
            if err == nil {
                err = stub.PutState(accountPrefix+renter.ID, renterBytes)
            } else {
                fmt.Println("Egads, something went wrong with Marshalling")
                return nil, errors.New("Egads, something went wrong with Marshalling")
            }
            
        } else {
            return nil, errors.New("Failed to add rent to owners") 
        }
    }


    return nil, nil

}

func (t *SimpleChaincode) issuePropertyToken(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

    /*      0
        json
        {
            "Name":  "name of the investment pool",
            "par": 0.00,
            "qty": 10,
            "discount": 7.5,
            "maturity": 30,
            "owners": [ // This one is not required
                {
                    "company": "company1",
                    "quantity": 5
                },
                {
                    "company": "company3",
                    "quantity": 3
                },
                {
                    "company": "company4",
                    "quantity": 2
                }
            ],              
            "issuer":"company2",
            "issueDate":"1456161763790"  (current time in milliseconds as a string)

        }
    */
    //need one arg
    if len(args) != 1 {
        fmt.Println("error invalid arguments")
        return nil, errors.New("Incorrect number of arguments. Expecting commercial paper record")
    }

    var cp PTY
    var err error
    var account Account

    var newstring = args[0]
    newstring = strings.Replace(args[0],"'","\"",-1)

    fmt.Println("Unmarshalling CP")
    err = json.Unmarshal([]byte(newstring), &cp)
    if err != nil {
        fmt.Println("error invalid paper issue")
        fmt.Println("error: ",err)
        return nil, errors.New("Invalid commercial paper issue")
    }

    fmt.Println("Hey guys, this is what we got:")
    fmt.Println("CP.name is   : ", cp.Name)
    fmt.Println("CP.Address is: ", cp.AdrStreet)
    fmt.Println("CP.Address is: ", cp.AdrCity)
    fmt.Println("CP.Address is: ", cp.AdrPostcode)
    fmt.Println("CP.Address is: ", cp.AdrState)
    cp.Status = "Pending"
    // Create string for hash

    stringHash := cp.AdrStreet+cp.AdrCity+cp.AdrPostcode+cp.AdrState

    cp.CUSIP, err = genHash(stringHash)

    fmt.Println("cusip is: ", cp.CUSIP)

    if cp.CUSIP == "" {
        fmt.Println("No CUSIP, returning error")
        return nil, errors.New("CUSIP cannot be blank")
    }
    fmt.Println("Getting state of - " + accountPrefix + cp.Issuer)
    accountBytes, err := stub.GetState(accountPrefix + cp.Issuer)
    if err != nil {
        fmt.Println("Error Getting state of - " + accountPrefix + cp.Issuer)
        return nil, errors.New("Error retrieving account " + cp.Issuer)
    }
    err = json.Unmarshal(accountBytes, &account)
    if err != nil {
        fmt.Println("Error Unmarshalling accountBytes")
        return nil, errors.New("Error retrieving account " + cp.Issuer)
    }
    
    //account.AssetsIds = append(account.AssetsIds, cp.CUSIP)

    var owner Owner
    owner.InvestorID = cp.Issuer
    owner.Quantity = cp.Qty

    cp.Owners = append(cp.Owners, owner)
    
    fmt.Println("Getting State on CP " + cp.CUSIP)
    cpRxBytes, err := stub.GetState(ptyPrefix+cp.CUSIP)
    if cpRxBytes == nil {
        fmt.Println("CUSIP does not exist, creating it")
        cpBytes, err := json.Marshal(&cp)
        if err != nil {
            fmt.Println("Error marshalling cp")
            return nil, errors.New("Error issuing commercial paper")
        }
        err = stub.PutState(ptyPrefix+cp.CUSIP, cpBytes)
        if err != nil {
            fmt.Println("Error issuing paper")
            return nil, errors.New("Error issuing commercial paper")
        }

        fmt.Println("Marshalling account bytes to write")
        accountBytesToWrite, err := json.Marshal(&account)
        if err != nil {
            fmt.Println("Error marshalling account")
            return nil, errors.New("Error issuing commercial paper")
        }
        err = stub.PutState(accountPrefix + cp.Issuer, accountBytesToWrite)
        if err != nil {
            fmt.Println("Error putting state on accountBytesToWrite")
            return nil, errors.New("Error issuing commercial paper")
        }
        
        
        // Update the paper keys by adding the new key
        fmt.Println("Getting Property Keys")
        keysBytes, err := stub.GetState("PtyKeys")
        if err != nil {
            fmt.Println("Error retrieving paper keys")
            return nil, errors.New("Error retrieving paper keys")
        }
        var keys []string
        err = json.Unmarshal(keysBytes, &keys)
        if err != nil {
            fmt.Println("Error unmarshel keys")
            return nil, errors.New("Error unmarshalling paper keys ")
        }
        
        fmt.Println("Appending the new key to Property Keys")
        foundKey := false
        for _, key := range keys {
            if key == ptyPrefix+cp.CUSIP {
                foundKey = true
            }
        }
        if foundKey == false {
            keys = append(keys, ptyPrefix+cp.CUSIP)
            keysBytesToWrite, err := json.Marshal(&keys)
            if err != nil {
                fmt.Println("Error marshalling keys")
                return nil, errors.New("Error marshalling the keys")
            }
            fmt.Println("Put state on Propert Keys")
            err = stub.PutState("PtyKeys", keysBytesToWrite)
            if err != nil {
                fmt.Println("Error writting keys back")
                return nil, errors.New("Error writing the keys back")
            }
        }
        fmt.Println("Issue Property Token %+v\n", cp)
        return nil, nil
    } else {
        fmt.Println("You can't tokenize an asset that already exists")
        return nil, errors.New("Can't tokenize asset that already exists")
    }
}

func (t *SimpleChaincode) setForSale(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    //   0
    // json
    // {
    //     CUSIP       string   `json:"cusip"`
    //     FromCompany string   `json:"fromCompany"`
    //     Quantity    int      `json:"quantity"`
    //     SellVal     float64  `json:"sellval"`
    // }

    //need one arg
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting commercial paper record")
    }
    
    var fs AddForSale

    fmt.Println("Unmarshalling ForSale")
    err := json.Unmarshal([]byte(strings.Replace(args[0],"'","\"",-1)), &fs)
    if err != nil {
        fmt.Println("Error Unmarshalling ForSale")
        return nil, errors.New("Invalid forsale issue")
    }

    fmt.Println("Getting State on CP " + fs.CUSIP)
    cpBytes, err := stub.GetState(ptyPrefix+fs.CUSIP)
    if err != nil {
        fmt.Println("CUSIP not found")
        return nil, errors.New("CUSIP not found " + fs.CUSIP)
    }

    var cp PTY
    fmt.Println("Unmarshalling CP " + fs.CUSIP)
    err = json.Unmarshal(cpBytes, &cp)
    if err != nil {
        fmt.Println("Error unmarshalling cp " + fs.CUSIP)
        return nil, errors.New("Error unmarshalling cp " + fs.CUSIP)
    }

    var fromCompany Account
    fmt.Println("Getting State on fromCompany " + fs.FromCompany)   
    fromCompanyBytes, err := stub.GetState(accountPrefix+fs.FromCompany)
    if err != nil {
        fmt.Println("Account not found " + fs.FromCompany)
        return nil, errors.New("Account not found " + fs.FromCompany)
    }

    fmt.Println("Unmarshalling FromCompany ")
    err = json.Unmarshal(fromCompanyBytes, &fromCompany)
    if err != nil {
        fmt.Println("Error unmarshalling account " + fs.FromCompany)
        return nil, errors.New("Error unmarshalling account " + fs.FromCompany)
    }

    // Check for all the possible errors
    ownerFound := false 
    quantity := 0
    for _, owner := range cp.Owners {
        if owner.InvestorID == fs.FromCompany {
            ownerFound = true
            quantity = owner.Quantity
        }
    }
    
    // If fromCompany doesn't own this paper
    if ownerFound == false {
        fmt.Println("The company " + fs.FromCompany + "doesn't own any of this paper")
        return nil, errors.New("The company " + fs.FromCompany + "doesn't own any of this paper")   
    } else {
        fmt.Println("The FromCompany does own this paper")
    }
    
    // If fromCompany doesn't own enough quantity of this paper
    if quantity < fs.Quantity {
        fmt.Println("The company " + fs.FromCompany + "doesn't own enough of this paper")       
        return nil, errors.New("The company " + fs.FromCompany + "doesn't own enough of this paper")            
    } else {
        fmt.Println("The FromCompany owns enough of this paper")
    }

    FromOwnerFound := false
    for key, owner := range cp.Owners {
        if owner.InvestorID == fs.FromCompany {
            fmt.Println("Reducing Quantity from the FromCompany")
            cp.Owners[key].Quantity -= fs.Quantity
//          owner.Quantity -= fs.Quantity
        }
    }
    for key, forsale := range cp.PT4Sale {
        if (forsale.InvestorID == fs.FromCompany) {
            FromOwnerFound = true
            fmt.Println("Found company in For Sale")
            cp.PT4Sale[key].Quantity += fs.Quantity
            cp.PT4Sale[key].SellVal = fs.SellVal
        }
    }
    
    if FromOwnerFound == false {
        var newOwner ForSale
        fmt.Println("As FromOwner was not found in ForSale, appending the owner to the CP")
        newOwner.Quantity = fs.Quantity
        newOwner.InvestorID = fs.FromCompany
        newOwner.SellVal = fs.SellVal
        cp.PT4Sale = append(cp.PT4Sale, newOwner)
    }

    // Write everything back
    // To Company
        
    // From company
    fromCompanyBytesToWrite, err := json.Marshal(&fromCompany)
    if err != nil {
        fmt.Println("Error marshalling the fromCompany")
        return nil, errors.New("Error marshalling the fromCompany")
    }
    fmt.Println("Put state on fromCompany")
    err = stub.PutState(accountPrefix+fs.FromCompany, fromCompanyBytesToWrite)
    if err != nil {
        fmt.Println("Error writing the fromCompany back")
        return nil, errors.New("Error writing the fromCompany back")
    }
    
    // cp
    cpBytesToWrite, err := json.Marshal(&cp)
    if err != nil {
        fmt.Println("Error marshalling the cp")
        return nil, errors.New("Error marshalling the cp")
    }
    fmt.Println("Put state on CP")
    err = stub.PutState(ptyPrefix+fs.CUSIP, cpBytesToWrite)
    if err != nil {
        fmt.Println("Error writing the cp back")
        return nil, errors.New("Error writing the cp back")
    }
    
    fmt.Println("Successfully completed Invoke")
    return nil, nil
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    //need one arg
    if len(args) < 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting ......")
    }

    if args[0] == "GetCompany" {
        fmt.Println("Getting the company")
        company, err := GetCompany(args[1], stub)
        if err != nil {
            fmt.Println("Error from getCompany")
            return nil, err
        } else {
            companyBytes, err1 := json.Marshal(&company)
            if err1 != nil {
                fmt.Println("Error marshalling the company")
                return nil, err1
            }   
            fmt.Println("All success, returning the company")
            return companyBytes, nil         
        }
    } else if args[0] == "GetAllPTYs" {
        fmt.Println("Getting all CPs")
        allCPs, err := GetAllPTYs(stub)
        if err != nil {
            fmt.Println("Error from GetAllPTYs")
            return nil, err
        } else {
            allCPsBytes, err1 := json.Marshal(&allCPs)
            if err1 != nil {
                fmt.Println("Error marshalling allptys")
                return nil, err1
            }   
            fmt.Println("All success, returning allptys")
            return allCPsBytes, nil      
        }
    } else if args[0] == "GetPTY" {
        fmt.Println("Getting all CPs")
        pty, err := GetPTY(args[1],stub)
        if err != nil {
            fmt.Println("Error from GetPTY")
            return nil, err
        } else {
            PtysBytes, err1 := json.Marshal(&pty)
            if err1 != nil {
                fmt.Println("Error marshalling ptys")
                return nil, err1
            }   
            fmt.Println("All success, returning ptys")
            return PtysBytes, nil      
        }
    } else {
        fmt.Println("I don't do shit!")
        fmt.Println("Generic Query call")
        bytes, err := stub.GetState(args[0])

        if err != nil {
            fmt.Println("Some error happenend")
            return nil, errors.New("Some Error happened")
        }

        fmt.Println("All success, returning from generic")
        return bytes, nil       
    }

    
    // if args[0] == "GetAllPTYs" {
    //     fmt.Println("Getting all CPs")
    //     allCPs, err := GetAllPTYs(stub)
    //     if err != nil {
    //         fmt.Println("Error from GetAllPTYs")
    //         return nil, err
    //     } else {
    //         allCPsBytes, err1 := json.Marshal(&allCPs)
    //         if err1 != nil {
    //             fmt.Println("Error marshalling allcps")
    //             return nil, err1
    //         }   
    //         fmt.Println("All success, returning allcps")
    //         return allCPsBytes, nil      
    //     }
    // }
    return nil, nil
}

func (t *SimpleChaincode) transferPaper(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    /*      0
        json
        {
              "CUSIP": "",
              "fromCompany":"",
              "toCompany":"",
              "quantity": 1
        }
    */
    //need one arg
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting commercial paper record")
    }
    
    var tr Transaction

    fmt.Println("Unmarshalling Transaction")
    err := json.Unmarshal([]byte(strings.Replace(args[0],"'","\"",-1)), &tr)
    if err != nil {
        fmt.Println("Error Unmarshalling Transaction")
        fmt.Println("err: ", err)
        return nil, errors.New("Invalid commercial paper issue")
    }

    fmt.Println("Getting State on CP " + tr.CUSIP)
    cpBytes, err := stub.GetState(ptyPrefix+tr.CUSIP)
    if err != nil {
        fmt.Println("CUSIP not found")
        return nil, errors.New("CUSIP not found " + tr.CUSIP)
    }

    var cp PTY
    fmt.Println("Unmarshalling CP " + tr.CUSIP)
    err = json.Unmarshal(cpBytes, &cp)
    if err != nil {
        fmt.Println("Error unmarshalling cp " + tr.CUSIP)
        return nil, errors.New("Error unmarshalling cp " + tr.CUSIP)
    }

    var fromCompany Account
    fmt.Println("Getting State on fromCompany " + tr.FromCompany)   
    fromCompanyBytes, err := stub.GetState(accountPrefix+tr.FromCompany)
    if err != nil {
        fmt.Println("Account not found " + tr.FromCompany)
        return nil, errors.New("Account not found " + tr.FromCompany)
    }

    fmt.Println("Unmarshalling FromCompany ")
    err = json.Unmarshal(fromCompanyBytes, &fromCompany)
    if err != nil {
        fmt.Println("Error unmarshalling account " + tr.FromCompany)
        return nil, errors.New("Error unmarshalling account " + tr.FromCompany)
    }

    var toCompany Account
    fmt.Println("Getting State on ToCompany " + tr.ToCompany)
    toCompanyBytes, err := stub.GetState(accountPrefix+tr.ToCompany)
    if err != nil {
        fmt.Println("Account not found " + tr.ToCompany)
        return nil, errors.New("Account not found " + tr.ToCompany)
    }

    fmt.Println("Unmarshalling tocompany")
    err = json.Unmarshal(toCompanyBytes, &toCompany)
    if err != nil {
        fmt.Println("Error unmarshalling account " + tr.ToCompany)
        return nil, errors.New("Error unmarshalling account " + tr.ToCompany)
    }

    // Check for all the possible errors
    ownerFound := false 
    quantity := 0
    price := 0.00
    for _, owner := range cp.PT4Sale {
        if owner.InvestorID == tr.FromCompany {
            ownerFound = true
            quantity = owner.Quantity
            price = owner.SellVal
        }
    }
    
    // If fromCompany doesn't own this paper
    if ownerFound == false {
        fmt.Println("The company " + tr.FromCompany + "doesn't own any of this paper")
        return nil, errors.New("The company " + tr.FromCompany + "doesn't own any of this paper")   
    } else {
        fmt.Println("The FromCompany does own this paper")
    }
    
    // If fromCompany doesn't own enough quantity of this paper
    if quantity < tr.Quantity {
        fmt.Println("The company " + tr.FromCompany + "doesn't own enough of this paper")       
        return nil, errors.New("The company " + tr.FromCompany + "doesn't own enough of this paper")            
    } else {
        fmt.Println("The FromCompany owns enough of this paper")
    }
    
    amountToBeTransferred := float64(tr.Quantity) * price
    
    // If toCompany doesn't have enough cash to buy the papers
    if toCompany.CashBalance < amountToBeTransferred {
        fmt.Println("The company " + tr.ToCompany + "doesn't have enough cash to purchase the papers")      
        return nil, errors.New("The company " + tr.ToCompany + "doesn't have enough cash to purchase the papers")   
    } else {
        fmt.Println("The ToCompany has enough money to be transferred for this paper")
    }

    // Checking to see if the shares are revofked
    if tr.FromCompany != tr.ToCompany {
        toCompany.CashBalance -= amountToBeTransferred
        fromCompany.CashBalance += amountToBeTransferred
    }

    toOwnerFound := false
    for key, owner := range cp.PT4Sale {
        if owner.InvestorID == tr.FromCompany {
            fmt.Println("Reducing Quantity from the FromCompany")
            cp.PT4Sale[key].Quantity -= tr.Quantity
//          owner.Quantity -= tr.Quantity
        }
        
    }

    for key, owner := range cp.Owners {
        if owner.InvestorID == tr.ToCompany {
            fmt.Println("Increasing Quantity from the ToCompany")
            toOwnerFound = true
            cp.Owners[key].Quantity += tr.Quantity
//          owner.Quantity += tr.Quantity
        }
    }
    
    if toOwnerFound == false {
        var newOwner Owner
        fmt.Println("As ToOwner was not found, appending the owner to the CP")
        newOwner.Quantity = tr.Quantity
        newOwner.InvestorID = tr.ToCompany
        cp.Owners = append(cp.Owners, newOwner)
    }
    
    //fromCompany.AssetsIds = append(fromCompany.AssetsIds, tr.CUSIP)

    // Write everything back
    // To Company
    toCompanyBytesToWrite, err := json.Marshal(&toCompany)
    if err != nil {
        fmt.Println("Error marshalling the toCompany")
        return nil, errors.New("Error marshalling the toCompany")
    }
    fmt.Println("Put state on toCompany")
    err = stub.PutState(accountPrefix+tr.ToCompany, toCompanyBytesToWrite)
    if err != nil {
        fmt.Println("Error writing the toCompany back")
        return nil, errors.New("Error writing the toCompany back")
    }
        
    // From company
    fromCompanyBytesToWrite, err := json.Marshal(&fromCompany)
    if err != nil {
        fmt.Println("Error marshalling the fromCompany")
        return nil, errors.New("Error marshalling the fromCompany")
    }
    fmt.Println("Put state on fromCompany")
    err = stub.PutState(accountPrefix+tr.FromCompany, fromCompanyBytesToWrite)
    if err != nil {
        fmt.Println("Error writing the fromCompany back")
        return nil, errors.New("Error writing the fromCompany back")
    }
    
    // cp
    cpBytesToWrite, err := json.Marshal(&cp)
    if err != nil {
        fmt.Println("Error marshalling the cp")
        return nil, errors.New("Error marshalling the cp")
    }
    fmt.Println("Put state on CP")
    err = stub.PutState(ptyPrefix+tr.CUSIP, cpBytesToWrite)
    if err != nil {
        fmt.Println("Error writing the cp back")
        return nil, errors.New("Error writing the cp back")
    }
    
    fmt.Println("Successfully completed Invoke")
    return nil, nil
}

func GetAllPTYs(stub *shim.ChaincodeStub) ([]PTY, error){
    
    var allCPs []PTY
    
    // Get list of all the keys
    keysBytes, err := stub.GetState("PtyKeys")
    if err != nil {
        fmt.Println("Error retrieving Property keys")
        return nil, errors.New("Error retrieving Property keys")
    }
    var keys []string
    err = json.Unmarshal(keysBytes, &keys)
    if err != nil {
        fmt.Println("Error unmarshalling Property keys")
        return nil, errors.New("Error unmarshalling Property keys")
    }

    // Get all the cps
    for _, value := range keys {
        cpBytes, err := stub.GetState(value)
        
        var cp PTY
        err = json.Unmarshal(cpBytes, &cp)
        if err != nil {
            fmt.Println("Error retrieving cp " + value)
            return nil, errors.New("Error retrieving cp " + value)
        }
        
        fmt.Println("Appending CP" + value)
        allCPs = append(allCPs, cp)
    }   
    
    return allCPs, nil
}

func GetPTY(cusip string, stub *shim.ChaincodeStub) (PTY, error){
    
    //
    cpBytes, err := stub.GetState(ptyPrefix+cusip)
    
    var cp PTY
    err = json.Unmarshal(cpBytes, &cp)
    if err != nil {
        fmt.Println("Error retrieving cp " + cusip)
        return cp, errors.New("Error retrieving cp " + cusip)
    }
    
    return cp, nil
}

func GetCompany(companyID string, stub *shim.ChaincodeStub) (Account, error){
    var company Account
    companyBytes, err := stub.GetState(accountPrefix+companyID)
    if err != nil {
        fmt.Println("Account not found " + companyID)
        return company, errors.New("Account not found " + companyID)
    }

    err = json.Unmarshal(companyBytes, &company)
    if err != nil {
        fmt.Println("Error unmarshalling account " + companyID + "\n err:" + err.Error())
        return company, errors.New("Error unmarshalling account " + companyID)
    }
    
    return company, nil
}

// Run callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will fsansfer X units from A to B upon invoke
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {


    fmt.Println("run is running " + function)
    return t.Invoke(stub, function, args)

}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    if function == "Init" {
        // Initialize the entities and their asset holdings
        return t.Init(stub,"init", args)
    } else if function == "issuePropertyToken" {
        // transaction makes payment of X units from A to B
        return t.issuePropertyToken(stub, args)
    } else if function == "createAccount" {
        // Deletes an entity from its state
        return t.createAccount(stub, args)
    } else if function == "createAccounts" {
        // Deletes an entity from its state
        return t.createAccounts(stub, args)
    } else if function == "setForSale" {
        // Deletes an entity from its state
        return t.setForSale(stub, args)
    } else if function == "transferPaper" {
        // Deletes an entity from its state
        fmt.Println("firing transferPaper")
        return t.transferPaper(stub, args)
    } else if function == "updateMktVal" {
        // Deletes an entity from its state
        return t.updateMktVal(stub, args)
    } else if function == "processRent" {
        // Deletes an entity from its state
        return t.processRent(stub, args)
    } else if function == "setRent" {
        // Deletes an entity from its state
        return t.setRent(stub, args)
    } else if function == "setRenters" {
       return t.setRenters(stub, args[0], args[1], args[2])
    }

    fmt.Println("Function"+ function +" was not found under invocation")
    return nil, errors.New("Received unknown function invocation")
}

func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Println("Error starting Simple chaincode: %s", err)
    }
}

var seventhDigit = map[int]string{
    1:  "A",
    2:  "B",
    3:  "C",
    4:  "D",
    5:  "E",
    6:  "F",
    7:  "G",
    8:  "H",
    9:  "J",
    10: "K",
    11: "L",
    12: "M",
    13: "N",
    14: "P",
    15: "Q",
    16: "R",
    17: "S",
    18: "T",
    19: "U",
    20: "V",
    21: "W",
    22: "X",
    23: "Y",
    24: "Z",
}

var eigthDigit = map[int]string{
    1:  "1",
    2:  "2",
    3:  "3",
    4:  "4",
    5:  "5",
    6:  "6",
    7:  "7",
    8:  "8",
    9:  "9",
    10: "A",
    11: "B",
    12: "C",
    13: "D",
    14: "E",
    15: "F",
    16: "G",
    17: "H",
    18: "J",
    19: "K",
    20: "L",
    21: "M",
    22: "N",
    23: "P",
    24: "Q",
    25: "R",
    26: "S",
    27: "T",
    28: "U",
    29: "V",
    30: "W",
    31: "X",
}

func (t *SimpleChaincode) calcRent(stub *shim.ChaincodeStub, args string) (float64) {

   CUSIP := args
   var p PTY
   var rentAmount float64

   byte, _ := stub.GetState(ptyPrefix+CUSIP)

   _ = json.Unmarshal(byte, &p)

   rentFloat := float64(len(p.Renters))
   rentAmount = p.Rent / rentFloat

   return rentAmount
}