import argparse
import json

from playwright.sync_api import sync_playwright


def scrape_with_playwright(url):
    """Scrapes the HTML content from the given URL using Playwright."""
    try:
        with sync_playwright() as p:
            browser = p.chromium.launch(headless=False)
            page = browser.new_page()
            page.goto(url)

            html_content = page.content()
            browser.close()
            return html_content, True, None
    except Exception as e:
        return None, False, e


def main():
    # Set up argument parsing
    parser = argparse.ArgumentParser(description="Scrape HTML content from a URL.")
    parser.add_argument("--url", type=str, required=True, help="The URL to scrape")

    args = parser.parse_args()

    # Scrape the provided URL
    html_content, success, error = scrape_with_playwright(args.url)

    data = {"html": html_content, "success": success, "error": error}

    print(json.dumps(data))

if __name__ == "__main__":
    main()
