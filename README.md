# Invoice

## Installation guide
1. clone the project
2. run mysql file or there is file db.txt open it then run only line one and line tow
3. if you did not run mysql file also do this step after runing first two lines run the project it creates tables for you then run last line of db.txt file
4. if you run db.sql it creates two users one is admin and one is user 
5. open project there is file .env update DB_USER to you mysql user  update DB_PASSWORD to user password update DB_NAME if you changes database name
6. if there is error with loadin packages there is file package.txt include all package commands to download


## Security
1. used jwt and refresh token
2. there is some ednpoints that related to user only admin have access to them like create, findAccounts and delete account

## TAX
1. you can update tax rate through .env file by change TAX_RATE
2. you can update threshold through .env file by change TAX_THRESHOLD



## REQUEST RESPONSE
1. it does not contain the response of delete requests because I need datas to test all project but you can also test delete requests i made ready for you
2. for product at first you can not put image but when product created you can update iamge product
3. all requests straight forward except invoiceUpdate which i expalin it in next section you can contact me for more information 

 ## Invoice Update Scenarios
 1. it has more than 9 scenarion
 2. update means invoiceUpdateLine  original means invoiceLineDto
 3. if update one has same length as original we need to compare update with original whether only quantity changes or itemId or both or may it is new one then check user deleted which invoice to delte it in db and create new ones
 4. check if it is for same customer other wise return money to the first customer subtract money from new customer also before everything check them balance
 5. same thing with more details happen when update greater than original or vice versa

## Test
1. first login with admin account userName is admin password is 1234
   
