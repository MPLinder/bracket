import json

from collections import Counter
from html.parser import HTMLParser


# For this to work you have to grab the json from the CBS players bracket page and then feed it to parse_bracket. The
# data is contained in a variabled called bootstrapBracketsData.
# Parse_bracket will not return teams picked to lose in the first round, and the counter value
# for each team will be 1 less than the number expected by the go scripts input json

def parse_picks(bracket):
    try:
        bracket = json.loads(bracket)
    except TypeError:
        pass

    all_teams = {}
    for team in bracket['game_and_pick_list']['teams']:
        all_teams[team['ceng_abbr']] = HTMLParser().unescape(team['name'])

    # A count of the CBS data will only net the number of games a player has picked a team to win
    # but the go script expects the number of games a team played, so we'll need to add 1 to every team.
    # We achieve this by starting the list out with one entry for every team.

    games_per_team = list(all_teams.values())
    for region in bracket['game_and_pick_list']['regions']:
        for round in region['rounds']:
            for game in round['games']:
                games_per_team.append(all_teams[game['user_pick']['pick']])

    return json.dumps(Counter(games_per_team))

def parse_actual(bracket):
    try:
        bracket = json.loads(bracket)
    except TypeError:
        pass

    all_teams = {}
    for team in bracket['game_and_pick_list']['teams']:
        all_teams[team['ceng_abbr']] = HTMLParser().unescape(team['name'])

    # A count of the CBS data will only net the number of games a player has picked a team to win
    # but the go script expects the number of games a team played, so we'll need to add 1 to every team.
    # We achieve this by starting the list out with one entry for every team.

    games_per_team = list(all_teams.values())
    for region in bracket['game_and_pick_list']['regions']:
        for round in region['rounds']:
            for game in round['games']:
                if game['winner_abbr'] != '':
                    games_per_team.append(all_teams[game['winner_abbr']])

    return json.dumps(Counter(games_per_team))
    
def parse_regions(bracket):
    try:
        bracket = json.loads(bracket)
    except TypeError:
        pass

    seed_map = {"top-left": 1, "bottom-left": 4, "top-right": 2, "bottom-right": 3}

    input_teams = bracket['game_and_pick_list']['teams']
    input_regions = bracket['game_and_pick_list']['regions']

    output_regions = []
    for region in input_regions:
        if region['name'] == 'finalfour':
            continue

        output_region = {}
        output_region['name'] = region['name']
        output_region['seed'] = seed_map[region['position']]

        output_region['teams'] = []
        for team in input_teams:
            if team['region_id'] == region['id']:
                team = {'name': HTMLParser().unescape(team['name']), 'seed': int(team['seed']), 'region': region['name']}
                output_region['teams'].append(team)

        output_regions.append(output_region)

    return json.dumps(output_regions)