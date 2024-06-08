# Web-dork👁
Webdork is a command-line application designed for targeted web dorking. It allows users to search specific domains and subdomains, as well as more general targets such as usernames or keywords, to uncover potential vulnerabilities and sensitive information. This tool can be particularly useful for cybersecurity professionals and penetration testers looking to gather intelligence about a target.

## How It Works

### 1. Command-line Flags

The application starts by parsing command-line flags to determine the mode of operation. The available flags are:

- `-d DOMAIN`: Specifies the target domain.
- `-s`: Indicates whether to include subdomains in the search.
- `-t TARGET`: Specifies a general target, such as a username or keyword.
- `-v`: Displays the version of the application.

### 2. Version Display

If the -v flag is used, the application prints the current version and exits. This is managed by the version() function which simply outputs the version number.

### 3. Dorking Mode

Depending on the flags provided, the application determines the mode of dorking:

Domain Dorking: If the -d flag is used, the application performs web dorking on the specified domain. If the -s flag is also used, it includes subdomains in the search.
Target Dorking: If the -t flag is used, the application performs web dorking on the general target.
The startDorking function manages the dorking process. It takes the target, a boolean indicating whether to include subdomains, and the type of target (domain or general).

### 4. Output Preparation

The application creates an output folder and an HTML file to store the results of the dorking process. This is done using the createOutputFile function, which ensures that the output directory exists and creates a new HTML file for the results.

### 5. Adding Dorks

The application reads dorks (search queries) from predefined text files:

dorks-domain.txt for domain dorking
dorks-target.txt for target dorking
The addDorks function reads the appropriate dork file, constructs Google search queries by combining each dork with the specified target, and writes these queries as clickable links into the HTML file.

### 6. Writing Results to HTML

Each constructed query is URL-encoded and added to the HTML file as a clickable link. This allows users to easily review the search results by opening the HTML file in a web browser. The HTML file is structured with a header and footer for better readability, managed by the addHTMLBanner and addHTMLFooter functions.

### 7. Opening in Browser

Once the HTML file is populated with the dork links, the application attempts to open the file in the user's default web browser. This is handled by the openInBrowser function, which uses platform-specific commands to ensure compatibility across different operating systems.

### 8. Web Scraping 

In addition to generating the HTML file, the application uses the colly library to perform web scraping on the Google search results. The performWebScraping function sets up a web scraper that visits the Google search URLs and extracts the actual search result links. These links are printed to the console for a quick overview.

### 9. Rate Limiting

To avoid triggering Google's anti-bot measures, the application includes a rate limit on web scraping requests. This is implemented in the performWebScraping function using the colly library's rate limiting feature, which ensures a delay between successive requests, mimicking human behavior.

## How to run the application

### Step 1 

```bash

go mod init dorking
```

### Step 2 

```bash

go run main.go
```


## Conclusion

Webdork is a powerful and flexible tool for web dorking, capable of targeting both specific domains and broader search queries. By automating the process of constructing and executing search queries, it simplifies the task of gathering intelligence on a target, making it an invaluable asset for security professionals.

## Acknowledgements

```bash
This app was made with 💖 by Hamza under the guidance of Sir Husnain.
```
