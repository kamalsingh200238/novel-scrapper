import re
from dataclasses import asdict, dataclass, field
from typing import Any, Dict, List

from botasaurus.browser import Driver, browser
from botasaurus.soupify import soupify
from chrome_extension_python import Extension
from ebooklib import epub

# Global variables for configuration
homepage_url = (
    "https://novelnext.dramanovels.io/nw/necromancer-of-the-shadows#tab-chapters-title"
)
start_chapter = 1044.0
end_chapter = 1158.0
author_name = "Kamal Singh"

# Dynamically construct the EPUB file name using the chapter range
epub_file_name = (
    f"necromancer_of_the_shadows_chapters_{int(start_chapter)}_to_{int(end_chapter)}.epub"
)


@dataclass
class HomepageResponse:
    html: str


@dataclass
class ChapterLink:
    url: str
    number: float


@dataclass
class ChapterContent:
    html: str
    number: float


@dataclass
class Ebook:
    title: str
    author: str
    chapters: List[ChapterContent] = field(default_factory=list)


def main():
    # Step 1: Scrape the homepage to get the list of chapter links
    homepage_response_dict = scrape_homepage(homepage_url)
    homepage_response = HomepageResponse(
        **homepage_response_dict
    )  # Convert back to dataclass
    chapter_links = extract_chapter_links(homepage_response.html)

    # Step 2: Filter the chapter links based on chapter range
    filtered_chapters = [
        link for link in chapter_links if start_chapter <= link.number <= end_chapter
    ]

    # Step 3: Scrape content of the filtered chapters
    chapter_contents_data = scrape_chapter_contents(filtered_chapters)

    # Convert chapter contents back to dataclass
    chapter_contents = [ChapterContent(**data) for data in chapter_contents_data]

    # Step 4: Create an Ebook object with the scraped chapters
    ebook_data = Ebook(
        title="Dragon Marked War God", author=author_name, chapters=chapter_contents
    )

    # Step 5: Generate the EPUB file
    create_epub(ebook_data)


def create_epub(ebook_data: Ebook):
    """
    Create an EPUB file from the given ebook data.
    """
    book = epub.EpubBook()

    # Set metadata
    book.set_identifier("id123456")
    book.set_title(ebook_data.title)
    book.set_language("en")
    book.add_author(ebook_data.author)

    # List to store chapter objects for TOC and spine
    chapter_items = []

    # Process each chapter and add it to the EPUB book
    for idx, chapter in enumerate(ebook_data.chapters):
        soup = soupify(chapter.html)

        scripts = soup.find_all("script")
        for script in scripts:
            script.extract()

        chapter_content = soup.find(id="chr-content")

        if chapter_content:
            chapter_title = f"Chapter {chapter.number}"
            epub_chapter = epub.EpubHtml(
                title=chapter_title, file_name=f"chapter_{idx + 1}.xhtml", lang="en"
            )
            epub_chapter.content = f"<h1>{chapter_title}</h1><p>{chapter_content.get_text()}</p>"

            # Add the chapter to the book and to the list for TOC and spine
            book.add_item(epub_chapter)
            chapter_items.append(epub_chapter)

    # Define the Table of Contents (TOC) and spine (reading order)
    book.toc = chapter_items  # Make sure it's a list, not a tuple
    book.spine = ["nav"] + chapter_items

    # Add necessary EPUB components (NCX, Nav, and CSS)
    book.add_item(epub.EpubNcx())
    book.add_item(epub.EpubNav())
    nav_css = epub.EpubItem(
        uid="style_nav",
        file_name="style/nav.css",
        media_type="text/css",
        content="body { font-family: Arial, sans-serif; color: black; }".encode(
            "utf-8"
        ),  # Encode to bytes
    )
    book.add_item(nav_css)

    # Write the EPUB file
    epub.write_epub(epub_file_name, book, {})
    print(f"EPUB file '{epub_file_name}' created successfully!")


@browser(
    block_images=True,
    headless=False,
)
def scrape_homepage(driver: Driver, url: str) -> Dict[str, Any]:
    """
    Scrape the homepage and return the HTML content as a dictionary.
    """
    driver.google_get(url)
    driver.wait_for_element("h3.title", 20)
    driver.scroll_to_bottom()
    driver.sleep(2)
    html = driver.page_html
    return asdict(HomepageResponse(html=html))  # Return as a dictionary


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
                ChapterLink(url=href, number=chapter_number)
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
)
def scrape_chapter_contents(driver: Driver, chapter: ChapterLink) -> Dict[str, Any]:
    """
    Scrape the content of the provided chapters and return a list of ChapterContent as dictionaries.
    """
    driver.google_get(chapter.url)
    driver.wait_for_element("a.novel-title", 20)
    html_content = driver.page_html
    return asdict(ChapterContent(html=html_content, number=chapter.number))


# Execute the main function to start the process
main()
