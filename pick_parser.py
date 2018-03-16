import HTMLParser
import json

from collections import Counter

# For this to work you have to grab the json from the CBS players bracket page and then feed it to parse_picks. The
# data is contained in a variabled called bootstrapBracketsData.
# Parse_picks will not return teams picked to lose in the first round, and the counter value
# for each team will be 1 less than the number expected by the go scripts input json

def parse_picks(picks_string):
    picks = json.loads(picks_string)

    all_teams = {}
    for team in picks['game_and_pick_list']['teams']:
        all_teams[team['ceng_abbr']] = HTMLParser.HTMLParser().unescape(team['name'])

    # A count of the CBS data will only net the number of games a player has picked a team to win
    # but the go script expects the number of games a team played, so we'll need to add 1 to every team.
    # We achieve this by starting the list out with one entry for every team.

    games_per_team = all_teams.values()
    for region in picks['game_and_pick_list']['regions']:
        for round in region['rounds']:
            for game in round['games']:
                games_per_team.append(all_teams[game['user_pick']['pick']])

    return json.dumps(Counter(games_per_team))

def parse_actual(picks_string):
    picks = json.loads(picks_string)

    all_teams = {}
    for team in picks['game_and_pick_list']['teams']:
        all_teams[team['ceng_abbr']] = HTMLParser.HTMLParser().unescape(team['name'])

    # Because I'm counting all games played here, I don't need to seed the return value the way parse_picks does

    games_per_team = all_teams.values()
    for region in picks['game_and_pick_list']['regions']:
        for round in region['rounds']:
            for game in round['games']:
                if game['winner_abbr'] != '':
                    games_per_team.append(all_teams[game['winner_abbr']])

    return json.dumps(Counter(games_per_team))
