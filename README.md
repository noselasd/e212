ABOUT
====
Web application to list E.212 MCC/MNC codes. 

This is a small/simple web application made for learning golang and
web development. It contains a bit og random features,

* Uses the go-macaron webapp framework
* Login/Logout of web application
* Authorization of access to certain URLs
* golang HTML templating
* Bootstrap4/JQuery/HTML5 experimentation
* HTTP post to do CRUD
* Json API to do CRUD
* sqlite3 database for data storage (No ORM, hand written SQL)

![Anonymous users](https://i.imgur.com/vhbBmaJ.png)

![Admin users](https://i.imgur.com/nodnCpc.png)


Build/Run
=========
1. Install golang (http://golang.org)
    When $GOPATH is not set, it refers to $HOME/go/

2. Install sqlite3 so you can run the sqlite3 command

3. clone the repository to your $GOPATH/go/src
   (default should just be to place it in $HOME/go/src, so the files are 
   within $HOME/go/src/e212/)

4. Install dependencies in vendor folder
   go mod vendor

5. cd $HOME/go/src/e212/

6. build the app:
    make

7. Add a user(admin with password admin)  and initialize the database:
    ./e212_cmd newuser admin admin@example.com admin

8. (Optional) Add E212 entries from ITU (Note, there can be some constraint violation
              due to duplicated entries from the ITU website)
    sqlite3 mccmnc.db < e212.sql

9. Run the app:
    ./e212

10. Visit the website at http://localhost:4000
    (Admin login can use the user:admin password:admin for the user created above)

Develop
=======
The app is developed in Linux in Visual Studio Code.
Install VSCode, and install the Go extensions.
Open the $GOPATH/go/e212/ folder.

TODO
====
* Move to github.com/noselasd/e212 and Use full import paths -
so project can also be fetched with go get
* Use ajax for form posts
* Make the About page
* Make documentation for the Json API


Author
======
Nils Olav Selåsdal - nos.utelsystems{at}gmail.com

