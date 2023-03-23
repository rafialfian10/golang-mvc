package handler

import (
	"12-database-relation-and-file-upload/connection"
	"12-database-relation-and-file-upload/model"
	"12-database-relation-and-file-upload/validation"
	"context"
	"strings"

	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux" // npm route
	"github.com/gorilla/sessions"

	// npm sessions
	"golang.org/x/crypto/bcrypt" // npm encript
)

// Function Route Home & Get Data
func HandleHome(w http.ResponseWriter, r *http.Request) { //ResponseWriter: untuk menampilkan data, Request: untuk menambahkan data
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // Header berfungsi untuk menampilkan data. Data yang ditamplikan "text-html" /"json" / dll

	tmpt, err := template.ParseFiles("views/index.html") // template.ParseFiles berfungsi memparsing file yang disisipkan sebagai parameter

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// Session. Tangkap session yang bernama SESSION_ID yang dikirm dari login
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	// Panggil struct Metadata untuk menampung data session
	var login model.MetaData

	// Masukkan data ke variabel login sesuai dengan kondisi login yang didapatkan dari session
	if session.Values["IsLogin"] != true {
		login.IsLogin = false
	} else {
		login.IsLogin = session.Values["IsLogin"].(bool)
		login.Name = session.Values["Name"].(string)
		login.Id = session.Values["Id"].(int)
	}

	// panggil method Flashes untuk mengkap data flash lalu buat array dan simpan ke var flashes
	fm := session.Flashes("message")
	var flashes []string

	// Jika fm lebih dari 0 / ada isinya maka simpan session
	if len(fm) > 0 {
		session.Save(r, w)

		// lakukan looping pada fm kemudian masukkan tiap data ke array flashes
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	// Kemudian join array dengan method strings.Join dan simpan kedalam variabel login (struct)
	login.FlashData = strings.Join(flashes, "")
	// End session

	// Query data dari database. Query mengembalikan 2 nilai
	dataProject, err := connection.Conn.Query(context.Background(), "SELECT tb_projects.id, project_name, start_date, end_date, description, technologies, images, tb_users.username FROM tb_projects LEFT JOIN tb_users ON tb_projects.author_id = tb_users.id ORDER BY id DESC")

	// dataProject, err := connection.Conn.Query(context.Background(), "SELECT id, project_name, start_date, end_date, description, technologies, images, 'Rafi' FROM tb_projects")
	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	// Panggil struct untuk menampung data dari database
	var result []model.Project

	// Sebelum data ditampilkan looping terlebih dahulu
	for dataProject.Next() {
		var data = model.Project{}

		// scan setiap data yang ada di struct dan database
		err := dataProject.Scan(&data.Id, &data.ProjectName, &data.StartDate, &data.EndDate, &data.Desc, &data.Tech, &data.Image, &data.Author)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, data)
	}

	// Buat map sebagai penampung data result
	data := map[string]interface{}{
		"Projects": result,
		"Login":    login,
	}

	// Kemudian tampilkan seluruh data dari database
	w.WriteHeader(http.StatusOK)
	tmpt.Execute(w, data)
}

// Function Route Project
func HandleProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/project.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	} else {
		tmpt.Execute(w, nil)
	}
}

// Function Route Contact
func HandleContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	temp, err := template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// Session. Tangkap session yang bernama SESSION_ID yang dikirm dari login
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	// Panggil struct Metadata untuk menampung data session
	var login model.MetaData

	// Masukkan data ke variabel login sesuai dengan kondisi login yang didapatkan dari session
	if session.Values["IsLogin"] != true {
		login.IsLogin = false
	} else {
		login.IsLogin = session.Values["IsLogin"].(bool)
		login.Name = session.Values["Name"].(string)
	}

	data := map[string]interface{}{
		"Login": login,
	}

	temp.Execute(w, data)
}

// Function Detail Project
func HandleDetailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/project-detail.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// Tangkap id dari blog
	id, _ := strconv.Atoi(mux.Vars(r)["id"]) // strconv.Atoi untuk konversi string ke int.  mux.Vars() berfungsi untuk menangkap id dan mengembalikan 2 nilai parameter result dan error

	// // Panggil struct untuk menampung data dari database
	var project model.Project

	// QueryRow(get 1 data) data dari database yang id didatabase sama dengan id yang ditangkap di URL
	row := connection.Conn.QueryRow(context.Background(), "SELECT tb_projects.id, project_name, start_date, end_date, description, technologies, images, tb_users.username FROM tb_projects LEFT JOIN tb_users ON tb_projects.author_id = tb_users.id WHERE tb_projects.id = $1", id)

	// Kemudia scan row
	err = row.Scan(&project.Id, &project.ProjectName, &project.StartDate, &project.EndDate, &project.Desc, &project.Tech, &project.Image, &project.Author)
	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	// Session
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	// Panggil struct Metadata untuk menampung data session
	var login model.MetaData

	// Masukkan data ke variabel login sesuai dengan kondisi login yang didapatkan dari session
	if session.Values["IsLogin"] != true {
		login.IsLogin = false
	} else {
		login.IsLogin = session.Values["IsLogin"].(bool)
	}

	// Buat map sebagai penampung data result
	dataProject := map[string]interface{}{
		"Project": project,
		"Login":   login,
	}

	tmpt.Execute(w, dataProject)
}

// Function Add Project
func HandleAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	temp, _ := template.ParseFiles("views/project.html")

	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	// Ambil value data input lalu tampung ke dalam variabel
	projectName := r.PostForm.Get("projectName")
	startDate := r.PostForm.Get("startDate")
	endDate := r.PostForm.Get("endDate")
	description := r.PostForm.Get("desc")
	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	// Buat array untuk menampung data checkbox
	var checkboxs []string

	// Jika didalam form checkboxs ada value-nya, maka append ke array checkboxs
	if r.FormValue("node") != "" {
		checkboxs = append(checkboxs, r.FormValue("node"))
	}
	if r.FormValue("angular") != "" {
		checkboxs = append(checkboxs, r.FormValue("angular"))
	}
	if r.FormValue("react") != "" {
		checkboxs = append(checkboxs, r.FormValue("react"))
	}
	if r.FormValue("typescript") != "" {
		checkboxs = append(checkboxs, r.FormValue("typescript"))
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	author := session.Values["Id"]

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(project_name, start_date, end_date, description, technologies, images, author_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", projectName, startDate, endDate, description, checkboxs, image, author)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	// Validation (Panggil function validation untuk melakukan validasi & pesan setelah data berhasil ditambahkan)
	validation := validation.NewValidation()
	project := model.Project{}

	data := make(map[string]interface{})
	vErrors := validation.Struct(project)

	// jika ada error tampilkan validasi, jika tidak tampilkan pesan
	if vErrors != nil {
		data["project"] = project
		// data["validation"] = vErrors
	} else {
		data["pesan"] = "Data project has been successfully added"
	}
	// End validaiton

	temp.Execute(w, data)
	// http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// Function Edit Project
func HandleEditProject(w http.ResponseWriter, r *http.Request) {
	// Jika request methodnya get, maka ambil data lamanya dan tampilkan didalm input
	if r.Method == http.MethodGet {
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		temp, err := template.ParseFiles("views/edit-project.html")

		if err != nil {
			w.Write([]byte("Message :" + err.Error()))
		}

		// Tangkap id
		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		// Panggil struct
		var project model.Project

		// QueryRow(get 1 data) data dari database yang id didatabase sama dengan id yang ditangkap di URL
		row := connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, description, technologies, images FROM tb_projects WHERE id = $1", id)

		// Scan hasil queryRow data dari database
		err = row.Scan(&project.Id, &project.ProjectName, &project.StartDate, &project.EndDate, &project.Desc, &project.Tech, &project.Image)
		if err != nil {
			w.Write([]byte("Message :" + err.Error()))
			return
		}

		// Buat map untuk menampung project yang telah discan. Kemudian didalam key "Project" buat map lagi lalu karena date akan diformat dulu
		dataProject := map[string]interface{}{
			"Project": map[string]interface{}{
				"Id":          project.Id,
				"ProjectName": project.ProjectName,
				"StartDate":   project.StartDate.Format("2006-01-02"),
				"EndDate":     project.EndDate.Format("2006-01-02"),
				"Desc":        project.Desc,
				"Tech":        project.Tech,
				"Image":       project.Image,
			},
		}

		// fmt.Println(dataProject)
		temp.Execute(w, dataProject)
		fmt.Println(dataProject)

	} else if r.Method == http.MethodPost {
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		temp, _ := template.ParseFiles("views/edit-project.html")

		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		// panggil struct kemudian simapn ke vaiabel new
		var new model.Project

		// Setiap value yang dipanggil dengan Form.Get akan masuk kedalam variabel new(struct Project)
		new.Id, _ = strconv.Atoi(r.Form.Get("id")) // id didapat dari input yang type-nya hidden
		new.ProjectName = r.FormValue("projectName")
		new.StartDate, _ = time.Parse("2006-01-02", r.FormValue("startDate")) // parsing value startDate ke time dengan method time.Parse
		new.EndDate, _ = time.Parse("2006-01-02", r.FormValue("endDate"))
		new.Desc = r.FormValue("desc")
		dataContext := r.Context().Value("dataFile")
		new.Image = dataContext.(string)

		// Jika didalam form checkboxs ada value-nya(node, angular, react, typescript), maka akan masuk kedalam variabel new
		if r.FormValue("node") != "" {
			new.Tech = append(new.Tech, r.FormValue("node"))
		}
		if r.FormValue("angular") != "" {
			new.Tech = append(new.Tech, r.FormValue("angular"))
		}
		if r.FormValue("react") != "" {
			new.Tech = append(new.Tech, r.FormValue("react"))
		}
		if r.FormValue("typescript") != "" {
			new.Tech = append(new.Tech, r.FormValue("typescript"))
		}

		fmt.Println("New Data :", new)

		// Kemudian panggil data dari database lalu UPDATE SET yang id dari database apakah sesuai dengan id di variabel new (dari URL)
		_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET project_name=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, images=$6 WHERE id=$7", new.ProjectName, new.StartDate, new.EndDate, new.Desc, new.Tech, new.Image, new.Id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Message :" + err.Error()))
			return
		}

		// Validation (Panggil function validation untuk melakukan validasi & pesan setelah data berhasil ditambahkan)
		validation := validation.NewValidation()
		project := model.Project{}

		data := make(map[string]interface{})
		vErrors := validation.Struct(project)

		if vErrors != nil {
			data["project"] = project
			// data["validation"] = vErrors
		} else {
			data["pesan"] = "Data project has been successfully updated"
		}

		temp.Execute(w, data)
		// End validaiton

		// http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

// Function Delete Project
func HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	// temp, _ := template.ParseFiles("views/index.html")

	// Tangkap id
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Kemudian delete data dimana id database sama dengan id yang ditangkap dari URL
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id = $1", id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
	}

	// Validation (Panggil function validation untuk melakukan validasi & pesan setelah data berhasil ditambahkan)
	validation := validation.NewValidation()
	project := model.Project{}

	data := make(map[string]interface{})
	vErrors := validation.Struct(project)

	if vErrors != nil {
		data["project"] = project
		// data["validation"] = vErrors
	} else {
		data["pesan"] = "Data project has been successfully deleted"
	}
	// End validaiton

	// temp.Execute(w, data)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// Function Route Register
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		temp, err := template.ParseFiles("views/register.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Message : " + err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		temp.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm()

		if err != nil {
			log.Fatal(err)
		}

		// Tangkap value
		username := r.PostForm.Get("username")
		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		// Encript password dengan method GenerateFromPassword kemudian masukkan password agar diencript
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

		// Panggil database dan masukkan value input ke database
		_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_users(username, email, password) VALUES($1, $2, $3)", username, email, passwordHash)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Message :" + err.Error()))
			return
		}

		// Buat session sebelum redirect ke halaman login
		store := sessions.NewCookieStore([]byte("SESSION_ID"))
		session, _ := store.Get(r, "SESSION_ID") // Panggil store.Get, store.Get mengembalikan 2 parameter(r, nama key session)

		session.Values["IsRegister"] = true
		session.AddFlash("Register successfully", "message") // AddFlash berfungsi untuk mengirimkan pesan. Par1 pesan, par2 message(key)
		session.Save(r, w)                                   // simpan session dengan method save dengan parameter request, response

		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
}

// Function Route Login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		temp, err := template.ParseFiles("views/login.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Message : " + err.Error()))
			return
		}

		// Panggil struct Metadata untuk menampung session
		var register model.MetaData

		// Buat session
		var store = sessions.NewCookieStore([]byte("SESSION_ID"))
		session, _ := store.Get(r, "SESSION_ID")

		// Kondisi jika IsLogin false maka return false, jika true maka return nilai IsLogin dan Name disimpan kedalam var login
		if session.Values["IsRegister"] != true {
			register.IsRegister = false
		} else {
			register.IsRegister = session.Values["IsRegister"].(bool)
		}

		fm := session.Flashes("message")
		var flashes []string

		if len(fm) > 0 {
			session.Save(r, w)

			for _, fl := range fm {
				flashes = append(flashes, fl.(string))
			}
		}
		register.FlashData = strings.Join(flashes, "")

		data := map[string]interface{}{
			"Login": register,
		}

		w.WriteHeader(http.StatusOK)
		temp.Execute(w, data)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm()

		if err != nil {
			log.Fatal(err)
		}

		// username := r.PostForm.Get("username")
		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		// Panggil struct user
		var user model.User

		err = connection.Conn.QueryRow(context.Background(), "SELECT id, username, email, password FROM tb_users WHERE email=$1", email).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Message :" + err.Error()))
			return
		}

		var store = sessions.NewCookieStore([]byte("SESSION_ID"))
		session, _ := store.Get(r, "SESSION_ID")

		// Panggil compareHashAndPassword untuk mencari tahu apakah password yang di inputkan itu sama / tidak dengan password di database yang telah dihash
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			temp, _ := template.ParseFiles("views/login.html")
			// w.Write([]byte("Message :" + err.Error()))
			session.Values["error"] = true
			session.AddFlash("Username and password is incorrect", "message")
			session.Save(r, w)

			var err model.MetaData

			if session.Values["error"] != true {
				err.IsLogout = false
			} else {
				err.IsLogout = session.Values["error"].(bool)
			}

			fm := session.Flashes("message")
			var flashes []string

			if len(fm) > 0 {
				session.Save(r, w)

				for _, fl := range fm {
					flashes = append(flashes, fl.(string))
				}
			}
			err.FlashData = strings.Join(flashes, "")
			data := map[string]interface{}{
				"Error": err,
			}
			temp.Execute(w, data)
			return
		}

		// // Buat session sebelum redirect ke halaman index

		session.Values["IsLogin"] = true
		session.Values["Name"] = user.Username
		session.Values["Id"] = user.Id
		session.Options.MaxAge = 10800                    // MaxAge berfungsi untuk menentukan jangka waktu maksimal pada cookie
		session.AddFlash("Login successfully", "message") // AddFlash berfungsi untuk mengrirmkan pesan (par1: isi pesan, par2: key)
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

// Function Logout
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Panggil SESSION_ID yang didapat dari cookie
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	// Set MaxAge pada cookie menjadi -1 agar cookie terhapus
	session.Values["IsLogout"] = true
	session.Options.MaxAge = -1
	session.AddFlash("Logout successfully", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
