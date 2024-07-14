# 🏫 UFABC Enrollment Filter

Webpage to filter UFABC enrollment. This project is separeted in two applications:

1. [Enrollment Filter](./enrollment-filter/) - Frontend interface to handle CSV files written using Vue.js.
2. [PDF Parser](./pdf-parser/) - A Golang code to download and parse a PDF to CSV.

## 📋 [Enrollment Filter](./enrollment-filter/)

This application allows users to filter and download enrollment data from a CSV file. Built with Vue.js, it offers a user-friendly interface to handle CSV files with ease.

### ⚙️ Prerequisites

Ensure you have the following installed:
- Node.js (>= 20.x)
- npm (>= 10.x)

1. Clone the repository:
    ```sh
    git clone https://github.com/Mewbi/ufabc-enrollment-filter.git
    cd ufabc-enrollment-filter/enrollment-filter
    ```

2. Install dependencies:
    ```sh
    npm install
    ```

### 🚀 Running the Application

#### 🖥️ Development Server

To start the development server, run:
```sh
npm run dev
```

The application will be available at `http://localhost:5173`.

#### 📦 Build for Production

To build the project for production, run:
```sh
npm run build
```

The output files will be generated in the `dist` directory.

### 🛠️ Inserting CSVs

Every new CSV must be moved to [public](./enrollment-filter/public/) directory and be referenced into `fileOptions` variable in the [HomePage](./enrollment-filter/src/views/HomePage.vue) file.

## 📄 [PDF Parser](./pdf-parser/)

The Go program downloads a PDF from a URL provided via the command line or via an endpoint API, parses table content from the PDF, and saves it as a CSV file or in a SQLite database. The table in the PDF should contain the columns "RA", "Código turma", and "Nome Disciplina" as the following [example](https://prograd.ufabc.edu.br/pdf/ajuste_2024_2_matriculas_deferidas.pdf):

<p align="center">
    <img src="./assets/example-pdf.png" height="300">
</p>

### ⚙️ Prerequisites

Before running the program, ensure you have Go installed on your system. You can download and install Go from [the official website](https://go.dev/doc/install).

### 🚀 Usage CLI

1. In your terminal navigate to the directory.

   ```sh
   cd pdf-parser/
   ```

2. Run the program with the URL of the PDF as a command-line argument:

   ```sh
   go run cmd/cli/main.go http://example.com/yourfile.pdf
   ```

   Replace `http://example.com/yourfile.pdf` with the actual URL of the PDF you want to process.

### 🚀 Usage API

1. In your terminal navigate to the directory.

   ```sh
   cd pdf-parser/
   ```

2. Start the API server with the following command:

   ```sh
   go run cmd/server/main.go
   ```


### 📂 Output

The program will download the PDF from the provided URL, parse the table data, and save the parsed data to a CSV file named `output.csv` in the same directory.

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

