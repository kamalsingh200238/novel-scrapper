import argparse
import re
from dataclasses import asdict, dataclass
from typing import Any, Dict, List

from botasaurus.browser import Driver, browser
from botasaurus.soupify import soupify
from botasaurus_driver.driver_utils import json
from chrome_extension_python import Extension


@dataclass
class HomepageResponse:
    html: str


@dataclass
class ChapterLink:
    number: float
    title: str
    url: str


@dataclass
class ChapterContent:
    number: float
    title: str
    content: str


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
    extensions=[
        Extension(
            "https://chromewebstore.google.com/detail/ublock-origin/cjpalhdlnbpafiamejdnhcphjbkeiagm"
        )
    ],
    output=None,
)
def scrape_homepage(driver: Driver, url: str) -> Dict[str, Any]:
    """
    Scrape the homepage and return the HTML content as a dictionary.
    """
    driver.google_get(url)
    driver.wait_for_element("h3.title", 60)
    driver.sleep(2)
    driver.scroll_to_bottom()
    driver.sleep(2)
    html = driver.page_html
    return asdict(HomepageResponse(html=html))


def extract_chapter_links(html: str) -> List[ChapterLink]:
    """
    Extract chapter links from the homepage HTML.
    """
    chapter_links = []
    soup = soupify(html)
    chapter_elements = soup.find_all("ul", class_="list-chapter")

    for chapter_list in chapter_elements:
        for i, link in enumerate(
            chapter_list.find_all("a")
        ):  # Enumerate to get the index
            href = link.get("href")
            title = link.get("title")
            # Updated regex to match both "Chapter 5(1)" and "Chapter 5.1"
            chapter_match = re.search(
                r"Chapter\s+(\d+)(?:\.(\d+)|\((\d+)\))?", title, re.IGNORECASE
            )

            if chapter_match:
                # Capture the main chapter number
                chapter_number = float(chapter_match.group(1))  # The main chapter part

                # Check if there's a decimal part or a parenthetical part
                if chapter_match.group(2):
                    chapter_number += float(
                        f"0.{chapter_match.group(2)}"
                    )  # Handle decimal part
                elif chapter_match.group(3):
                    chapter_number += float(
                        f"0.{chapter_match.group(3)}"
                    )  # Handle parentheses part
            else:
                # If regex does not match, set chapter_number to the loop index i
                chapter_number = float(i)  # Use the loop index i as the default value

            # Append the chapter link with the calculated chapter number
            chapter_links.append(
                ChapterLink(number=chapter_number, title=title, url=href)
            )  # Append as ChapterLink instance

    return chapter_links


@browser(
    block_images=True,
    headless=False,
    reuse_driver=True,
    extensions=[
        Extension(
            "https://chromewebstore.google.com/detail/ublock-origin/cjpalhdlnbpafiamejdnhcphjbkeiagm"
        )
    ],
    output=None,
)
def scrape_chapter_contents(driver: Driver, chapter: ChapterLink) -> Dict[str, Any]:
    """
    Scrape the content of the provided chapters and return a list of ChapterContent as dictionaries.
    """
    driver.google_get(chapter.url)
    driver.wait_for_element("a.novel-title", 60)

    html = driver.page_html
    soup = soupify(html)

    scripts = soup.find_all("script")
    for script in scripts:
        script.extract()

    chapter_html = soup.find(id="chr-content")
    chapter_content = chapter_html.get_text()

    return asdict(
        ChapterContent(
            content=chapter_content, number=chapter.number, title=chapter.title
        )
    )


def main():
    args = parse_arguments()

    try:
        homepage_response_dict = scrape_homepage(args.url)
        homepage_response = HomepageResponse(**homepage_response_dict)
        chapter_links = extract_chapter_links(homepage_response.html)
        filtered_chapters = [
            link for link in chapter_links if args.start <= link.number <= args.end
        ]
        chapter_contents_dict = scrape_chapter_contents(filtered_chapters)
        resp = {
            "data": chapter_contents_dict,
            "error": "",
        }
        print(json.dumps(resp))
    except Exception as e:
        resp = {
            "data": "",
            "error": e,
        }
        print(json.dumps(resp))


if __name__ == "__main__":
    main()
