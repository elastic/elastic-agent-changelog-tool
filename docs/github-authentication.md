# GitHub authentication

Some commands require access to the GitHub APIs.

The tool uses the GitHub token to authorize user's call to API.  
The token can be stored in the `~/.elastic/github.token` file or passed via the `GITHUB_TOKEN` environment variable.

Here are the instructions on how to create your own personal access token (PAT): [GitHub docs - Creating a personal access token][1].

Make sure you have enabled the following scopes:
* `public_repo` — to open pull requests on GitHub repositories.
* `read:user` and `user:email` — to read your user profile information from GitHub.

**NOTE**: Elastic GitHub organization uses SAML authentication, due to this PAT cannot be used for Elastic related repositories **unless authorized**.
After creating or modifying your personal access token, authorize the token for use of the Elastic organization: [GitHub docs - Authorizing a personal access token for use with SAML single sign-on][2].

[1]: https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token
[2]: https://docs.github.com/en/github/authenticating-to-github/authenticating-with-saml-single-sign-on/authorizing-a-personal-access-token-for-use-with-saml-single-sign-on

