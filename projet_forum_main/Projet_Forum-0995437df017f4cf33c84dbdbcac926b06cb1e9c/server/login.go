package connexion

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	cryptage "./crypt"
	_ "github.com/mattn/go-sqlite3"
)

func display_db(db *sql.DB) {
	query := "SELECT * FROM Utilisateur"
	result, err := db.Query(query)
	if err != nil {
		println("utilisateur n'existe pas")
	}
	var PASSWORD string
	var MAIL string
	var Nom string
	var PRENOM string
	var ID_user int
	var ADDRESSE string
	var Date string
	for result.Next() {
		result.Scan(&ID_user, &Nom, &PRENOM, &MAIL, &PASSWORD, &ADDRESSE, &Date)
	}
}
func select_password(db *sql.DB, address string) string {
	query := "SELECT * FROM Utilisateur WHERE MAIL='" + address + "'"
	result, err := db.Query(query)
	if err != nil {
		println("utilisateur n'existe pas")
	}
	var PASSWORD string
	var MAIL string
	var Nom string
	var PRENOM string
	var ID_user int
	var ADDRESSE string
	var Date string
	for result.Next() {
		result.Scan(&ID_user, &Nom, &PRENOM, &MAIL, &PASSWORD, &ADDRESSE, &Date)
		return PASSWORD
	}
	return "Utilisateur n'existe pas dans la base"
}
func initdatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func renderTemplate_creation(w http.ResponseWriter, r *http.Request) {

	parsedTemplate, _ := template.ParseFiles("./template/creation_compte.html")
	//Call to ParseForm makes form fields available.
	err := r.ParseForm()
	if err != nil {
		print("Error\n")
		// Handle error here via logging and then return
	}

	// the values ​​of the form
	Prenom := r.PostFormValue("first_name")
	Nom := r.PostFormValue("last_name")
	User_name := r.PostFormValue("User_name")
	MDP := r.PostFormValue("MDP")
	Mail := r.PostFormValue("Mail")
	Date := r.PostFormValue("date_naissance")
	println(Date)
	// open the database (we create it if it does not exist)
	database, err := sql.Open("sqlite3", "./Forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// insert values ​​into the database with the INSERT INTO request
	statement, _ := database.Prepare("INSERT INTO Utilisateur (Nom, PRENOM,MAIL,PASSWORD,User_name,Birth_Date) VALUES (?, ?,?,?,?,?)")
	MDP_Hash, _ := cryptage.HashPassword(MDP)
	// we insert in the database if the values ​​are not empty
	if Nom != "" && Prenom != "" && Mail != "" && MDP != "" && User_name != "" {
		statement.Exec(Nom, Prenom, Mail, MDP_Hash, User_name, Date)
	}
	err_tmpl := parsedTemplate.Execute(w, nil)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
func renderTemplate_login(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("./template/login.html")
	database := initdatabase("./Forum.db")
	display_db(database)
	//Call to ParseForm makes form fields available.
	err := r.ParseForm()
	if err != nil {
		print("Error\n")
		// Handle error here via logging and then return
	}
	MDP := r.PostFormValue("MDP")
	Mail := r.PostFormValue("Mail")
	password := select_password(database, Mail)
	println(password)
	if cryptage.Verif(MDP, password) {
		http.SetCookie(w, &http.Cookie{
			Name:  "logged-in",
			Value: "1",
			Path:  "/",
		})
		http.Redirect(w, r, "/Accueil.html", http.StatusFound)
		println("tout est bon")
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:  "logged-in",
			Value: "0",
			Path:  "/",
		})
		println("faux mot de passe")
		// http.Redirect(w, r, "/login.html", http.StatusFound)
	}
	err_tmpl := parsedTemplate.Execute(w, nil)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
func renderTemplate_post(w http.ResponseWriter, r *http.Request) {

	parsedTemplate, _ := template.ParseFiles("./template/creation_post.html")
	//Call to ParseForm makes form fields available.
	err := r.ParseForm()
	if err != nil {
		print("Error\n")
		// Handle error here via logging and then return
	}
	database, err := sql.Open("sqlite3", "./Forum.db")
	if err != nil {
		log.Fatal(err)
	}

	titre := r.FormValue("Titre")
	message := r.FormValue("Message")
	informatique := r.FormValue("Informatique")
	jv := r.FormValue("Jv")
	art := r.FormValue("Art")
	sport := r.FormValue("Sport")
	ht := r.FormValue("Ht")
	cuisine := r.FormValue("Cuisine")

	statements, _ := database.Prepare(`INSERT INTO post (titre, message, informatique, jv, art, sport, ht, cuisine) VALUES (?,?,?,?,?,?,?,?)`)

	if titre != "" && message != "" {
		statements.Exec(titre, message, informatique, jv, art, sport, ht, cuisine)
	}
	err_tmpl := parsedTemplate.Execute(w, nil)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
func Post(data_info map[string]interface{}, id string) {
	var post [][]interface{}

	database, err := sql.Open("sqlite3", "./Forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	//range over database
	var query string
	if id == "allposts" {
		query = "SELECT ID_user, titre, message, informatique, jv, art, sport, ht, cuisine FROM post"
	} else {
		query = "SELECT ID_user, titre, message, informatique, jv, art, sport, ht, cuisine FROM post WHERE ID_user = " + id
	}
	rows, err := database.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		aPost := []interface{}{"", "", "", "", "", "", "", "", ""}
		rows.Scan(&aPost[0], &aPost[1], &aPost[2], &aPost[3], &aPost[4], &aPost[5], &aPost[6], &aPost[7], &aPost[8])
		post = append(post, aPost)
	}
	data_info["allposts"] = post
}
func renderTemplate_accueil(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("logged-in")
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	if r.URL.Path == "/logout.html" {
		http.SetCookie(w, &http.Cookie{
			Name:  "logged-in",
			Value: "0",
			Path:  "/",
		})
		http.Redirect(w, r, "/Accueil.html", http.StatusFound)
	}
	println(c.Value)
	parsedTemplate, _ := template.ParseFiles("./template/Accueil.html")
	err_tmpl := parsedTemplate.Execute(w, nil)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
func pagepost(w http.ResponseWriter, r *http.Request) {
	data_info := make(map[string]interface{})
	Post(data_info, "allposts")

	parsedTemplate, _ := template.ParseFiles("./template/pagepost.html")
	err_tmpl := parsedTemplate.Execute(w, data_info)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
func Comments(data_info map[string]interface{}, ID_user string) {
	// Comments
	var comments [][1]string
	database, err := sql.Open("sqlite3", "./Forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	//range over database
	rows_comment, err := database.Query("SELECT commente FROM reponse WHERE ID_post = ?", ID_user)
	if err != nil {
		log.Fatal(err)
	}
	defer rows_comment.Close()
	for rows_comment.Next() {
		var tmp [1]string
		err := rows_comment.Scan(&tmp[0])
		if err != nil {
			log.Fatal(err)
		}
		comments = append(comments, tmp)
	}
	err = rows_comment.Err()
	if err != nil {
		log.Fatal(err)
	}
	data_info["comments"] = comments
}
func renderTemplate_reponse(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	data_info := make(map[string]interface{})
	Post(data_info, id)

	parsedTemplate, _ := template.ParseFiles("./template/reponse_post.html")
	//Call to ParseForm makes form fields available.
	err := r.ParseForm()
	if err != nil {
		print("Error\n")
		// Handle error here via logging and then return
	}

	database, err := sql.Open("sqlite3", "./Forum.db")
	if err != nil {
		log.Fatal(err)
	}

	commente := r.PostFormValue("Commentaire")

	statements, _ := database.Prepare(`INSERT INTO reponse (ID_post,commente) VALUES (?,?)`)
	if commente != "" {
		statements.Exec(id, commente)
	}
	Comments(data_info, id)

	err_tmpl := parsedTemplate.Execute(w, data_info)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
func Info(w http.ResponseWriter, r *http.Request) {
	data_info := make(map[string]interface{})
	Post(data_info, "allposts")

	parsedTemplate, _ := template.ParseFiles("./template/informatique_post.html")
	err_tmpl := parsedTemplate.Execute(w, data_info)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
func JeuxVideo(w http.ResponseWriter, r *http.Request) {
	data_info := make(map[string]interface{})
	Post(data_info, "allposts")

	parsedTemplate, _ := template.ParseFiles("./template/jeuxvideo_post.html")
	err_tmpl := parsedTemplate.Execute(w, data_info)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
func Art(w http.ResponseWriter, r *http.Request) {
	data_info := make(map[string]interface{})
	Post(data_info, "allposts")

	parsedTemplate, _ := template.ParseFiles("./template/art_post.html")
	err_tmpl := parsedTemplate.Execute(w, data_info)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
func Sport(w http.ResponseWriter, r *http.Request) {
	data_info := make(map[string]interface{})
	Post(data_info, "allposts")

	parsedTemplate, _ := template.ParseFiles("./template/sport_post.html")
	err_tmpl := parsedTemplate.Execute(w, data_info)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
func HighTech(w http.ResponseWriter, r *http.Request) {
	data_info := make(map[string]interface{})
	Post(data_info, "allposts")

	parsedTemplate, _ := template.ParseFiles("./template/hightech_post.html")
	err_tmpl := parsedTemplate.Execute(w, data_info)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
func Cuisine(w http.ResponseWriter, r *http.Request) {
	data_info := make(map[string]interface{})
	Post(data_info, "allposts")

	parsedTemplate, _ := template.ParseFiles("./template/cuisine_post.html")
	err_tmpl := parsedTemplate.Execute(w, data_info)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
func Login() {
	db := initdatabase("Forum.db")
	defer db.Close()
	http.HandleFunc("/Accueil.html", renderTemplate_accueil)
	http.HandleFunc("/creation_compte.html", renderTemplate_creation)
	http.HandleFunc("/login.html", renderTemplate_login)
	http.HandleFunc("/creation_post.html", renderTemplate_post)
	http.HandleFunc("/pagepost.html", pagepost)
	http.HandleFunc("/reponse_post.html", renderTemplate_reponse)
	http.HandleFunc("/informatique_post.html", Info)
	http.HandleFunc("/jeuxvideo_post.html", JeuxVideo)
	http.HandleFunc("/art_post.html", Art)
	http.HandleFunc("/sport_post.html", Sport)
	http.HandleFunc("/hightech_post.html", HighTech)
	http.HandleFunc("/cuisine_post.html", Cuisine)
	fs := http.FileServer(http.Dir("./assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
}
