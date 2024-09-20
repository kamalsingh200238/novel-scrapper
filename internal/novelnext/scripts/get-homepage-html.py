import argparse
import time
import json

from playwright.sync_api import sync_playwright


def scrape_with_playwright(url):
    """Scrapes the HTML content from the given URL using Playwright."""
    try:
        with sync_playwright() as p:
            browser = p.chromium.launch(headless=True)
            page = browser.new_page()
            page.goto(url)

            # Scroll to the bottom of the page to load all dynamic content
            last_height = page.evaluate("document.body.scrollHeight")
            while True:
                page.evaluate("window.scrollTo(0, document.body.scrollHeight)")
                time.sleep(2)
                new_height = page.evaluate("document.body.scrollHeight")
                if new_height == last_height:
                    break
                last_height = new_height

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
