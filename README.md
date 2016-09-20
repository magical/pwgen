Pwgen is a minimalist online password generator.

Live version at <https://turnipmints.mooo.com/password>.

### Usage

Go to <https://turnipmints.mooo.com/password>.
Pick four or five of your favorite words from the list.
That's your new password!

### History

I used to use a python script and a custom word list to generate
memorable passwords, but they got left behind on an old computer at some point,
forcing me to use diceware and shuf(1), or resort to online generators.
But I don't really like diceware, and online generators are usually overly complex.
It was time to write my own. Fortunately, around this time i discovered that
the EFF had published a new and improved word list. Thus pwgen was born.

### Installation

    $ go install github.com/magical/pwgen
    $ ./pwgen -bind :8080 &
    $ firefox http://localhost:8080/password

### License

The code is licensed under the WTFPL version 2.0. No warranty.

The word lists are from the EFF and are licensed under the
Creative Commons Attribution License 3.0, per <https://www.eff.org/copyright>.
