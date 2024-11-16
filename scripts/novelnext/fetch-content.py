import argparse
from dataclasses import asdict, dataclass
from typing import Any, Dict

from botasaurus.browser import Driver, browser


@dataclass
class HomepageResponse:
    html: str


def parse_arguments():
    parser = argparse.ArgumentParser(
        description="Script with required command line arguments"
    )
    parser.add_argument(
        "--url", type=str, required=True, help="The URL to novel's homepage"
    )
    parser.add_argument(
        "--start", type=float, required=True, help="Start chapter number (float)"
    )
    parser.add_argument(
        "--end", type=float, required=True, help="End chapter number (float)"
    )
    return parser.parse_args()


@browser(
    block_images=True,
    headless=False,
)
def scrape_homepage(driver: Driver, url: str) -> Dict[str, Any]:
    """
    Scrape the homepage and return the HTML content as a dictionary.
    """
    driver.google_get(url)
    driver.wait_for_element("h3.title", 100)
    driver.sleep(2)
    driver.scroll_to_bottom()
    driver.sleep(2)
    html = driver.page_html
    return asdict(HomepageResponse(html=html))


def main():
    args = parse_arguments()

    print(f"URL: {args.url}")
    print(f"Start: {args.start}")
    print(f"End: {args.end}")

    # get all the links


if __name__ == "__main__":
    main()
