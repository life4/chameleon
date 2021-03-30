import re
import sys
import tokenize
import typing
from pathlib import Path

REX_WORD = re.compile(r'[A-Za-z]+')
# https://stackoverflow.com/a/1176023/8704691
REX1 = re.compile(r'(.)([A-Z][a-z]+)')
REX2 = re.compile(r'([a-z0-9])([A-Z])')


def get_words(text: str) -> typing.Set[str]:
    text = REX1.sub(r'\1 \2', text)
    text = REX2.sub(r'\1 \2', text).lower()
    text = text.replace('_', ' ')
    return set(text.split())


def get_redundant_comments(tokens: typing.Iterable[tokenize.TokenInfo]):
    comment = None
    code_words: typing.Set[str] = set()
    for token in tokens:
        if token.type == tokenize.NL:
            continue

        if token.type == tokenize.COMMENT:
            if comment is None:
                comment = token
            else:
                comment = None
            continue
        if comment is None:
            continue

        if token.start[0] == comment.start[0] + 1:
            code_words.update(get_words(token.string))
            continue

        comment_words = set(REX_WORD.findall(comment.string))
        if comment_words and not comment_words - code_words:
            yield comment
        code_words = set()
        comment = None


def process_file(path: Path):
    with path.open('rb') as stream:
        tokens = tokenize.tokenize(stream.readline)
        for comment in get_redundant_comments(tokens):
            print(f'{path}:{comment.start[0]} {comment.string.strip()}')


def get_paths(path: Path):
    if path.is_dir():
        for p in path.iterdir():
            yield from get_paths(p)
    elif path.suffix == '.py':
        yield path


if __name__ == '__main__':
    for path in sys.argv[1:]:
        for p in get_paths(Path(path)):
            # process file
            process_file(p)
