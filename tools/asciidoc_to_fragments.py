# Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
# or more contributor license agreements. Licensed under the Elastic License 2.0;
# you may not use this file except in compliance with the Elastic License 2.0.

import argparse
import requests

from os import makedirs
from os.path import expanduser
from datetime import datetime


api_url = "https://api.github.com/repos/"
github_token_location = "/.elastic/github.token"

fragments_path = "changelog/fragments/"
fragments_counter = 0

repo_dict = {
    "agent": "https://github.com/elastic/elastic-agent",
    "beats": "https://github.com/elastic/beats"
}
default_repolink = "https://github.com/elastic/elastic-agent"

kind_dict = {
    "Breaking changes": "breaking-change",
    "Bugfixes": "bug-fix",
    "New features": "feature",
}
kind_token = "===="
field_token = "-"

def write_fragment(title, fragment_timestamp, fragment_dict):
    if not fragment_timestamp:
        fragment_timestamp = str(1000000000 + fragments_counter)

    path = "".join([fragments_path,
                    fragment_timestamp,
                    "-",
                    title,
                    ".yaml"])

    with open(path, 'w+') as f:
        for k, v in fragment_dict.items():
            f.write(f"{k}: {v}\n")

def get_event_timestamp(repository, event, number):
    token_path = ''.join([expanduser("~"), github_token_location])
    with open(token_path, 'r') as f:
        token = f.read().rstrip()

    owner, repo = repository.split("/")[-2:]
    event_url = f"{owner}/{repo}/{event}/{number}"
    url = f"{api_url}{event_url}"
    headers = {"Accept": "application/vnd.github+json", "Authorization": f"Bearer {token}"}
    response = requests.get(url, headers=headers)
    data = response.json()

    if response.status_code == 404:
        return "not_found"
    elif data["closed_at"] is None:
        return "event_open"
    else:
        date = datetime.fromisoformat(data["closed_at"].replace('Z', '+00:00'))
        return str(int(datetime.timestamp(date)))

def parse_line(line, kind):
    global fragments_counter
    fragments_counter += 1

    summary, *entries = line.split(" {")
    if len(entries) == 0:
        print(f"Warning: {line} -> no PR/issue fields\n")

    fragment_dict = {"kind": kind}
    fragment_dict["summary"] = summary.lstrip(field_token).strip()
    fragment_dict["summary"] = fragment_dict["summary"].replace(":", "")

    title = fragment_dict["summary"]
    title = title.replace(" ", "-")
    title = title.replace("/", "|")
    title = title.rstrip(".")

    pr_repo, issue_repo, fragment_timestamp = "", "", ""

    for entry in entries:
        number = entry[entry.find("[")+1:entry.find("]")]
        number = ''.join(filter(lambda n: n.isdigit(), number))
        entry_data = entry.split("}")[0]

        try:
            fragment_field, repo = entry_data.split("-")
            repo_link = repo_dict[repo]
        except ValueError:
            fragment_field, repo_link = entry_data, default_repolink

        if fragment_field in fragment_dict.keys():
            print(f"Skipping {line} -> multiple PRs/issues found\n")
            return

        if fragment_field == "pull":
            fragment_dict["pr"] = ''.join([repo_link, '/pull/', number])
            pr_number, pr_repo = number, repo_link
        elif fragment_field == "issue":
            fragment_dict["issue"] = ''.join([repo_link, '/issues/', number])
            issue_number, issue_repo = number, repo_link
    
    if pr_repo:
        fragment_timestamp = get_event_timestamp(pr_repo, "pulls", pr_number)
    elif issue_repo:
        fragment_timestamp = get_event_timestamp(issue_repo, "issues", issue_number)

    if fragment_timestamp == "not_found":
        print(f"Skipping {line} -> no response from Github API\n")
        return
    if fragment_timestamp == "event_open":
        print(f"Skipping {line} -> PR/issue still open!\n")
        return

    if issue_repo != pr_repo and pr_repo:
        try:
            del fragment_dict["issue"]
            print(f"Warning: {line} -> issue info lost due to multiple repositories\n")
        except KeyError:
            pass

    write_fragment(title, fragment_timestamp, fragment_dict)

def iterate_lines(f, kind='', skip=True):
    line = next(f, None)
    if line is None:
        return

    if line.startswith(kind_token):
        iterate_lines(f, kind_dict[line.lstrip(kind_token).strip()], skip=False)

    elif line.isspace():
        iterate_lines(f, kind, skip)

    elif line.startswith(field_token) and skip is False:
        parse_line(line, kind)

    else:
        iterate_lines(f, kind, skip=True)

    iterate_lines(f, kind, skip)

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--path", help="Changelog file path", required=True)
    parser.add_argument("--workdir", help="Working directory path")
    parser.add_argument("--repo", help="Repository name")
    args = parser.parse_args()

    if args.workdir:
        args.path = ''.join([args.workdir, '/', args.path])
        fragments_path = ''.join([args.workdir, '/', fragments_path])

    if args.repo:
        default_repolink = repo_dict[args.repo]

    try:
        makedirs(fragments_path)
    except FileExistsError as e:
        pass

    print("Skipped entries should be manually created and warnings should be checked")
    with open(args.path, 'r') as f:
        iterate_lines(f)
