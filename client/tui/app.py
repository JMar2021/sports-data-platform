from textual.app import App
from textual.widgets import Static, DataTable, Header
from textual.worker import Worker, WorkerState
from textual import work
import requests

class SportsDataApp(App):
    def compose(self):
        yield Header()
        yield Static("Loading...")
        yield DataTable()

    def on_mount(self):
        table = self.query_one(DataTable)

        table.add_columns(
            "Away",
            "Score",
            "Home",
            "Score",
            "Status",
            "Venue",
        )

        self.refresh_games()
        self.set_interval(10, self.refresh_games)

    @work(thread=True, group="refresh", exclusive=True)
    def refresh_games(self):
        return get_games()

    def on_worker_state_changed(self, event: Worker.StateChanged):
        if event.state == WorkerState.SUCCESS:
            games = event.worker.result

            table = self.query_one(DataTable)
            table.clear()

            for game in games:
                table.add_row(
                    game["AwayTeam"],
                    str(game["AwayScore"]),
                    game["HomeTeam"],
                    str(game["HomeScore"]),
                    game["Status"],
                    game["Venue"],
                )
            status = self.query_one(Static)
            status.update("Last updated successfully")

        elif event.state == WorkerState.ERROR:
            stauts.update(f"Refresh failed: {event.worker.error}")

def get_games():
    response = requests.get("http://localhost:8080/games", timeout = 10)
    response.raise_for_status()
    return response.json()

if __name__ == "__main__":
    app = SportsDataApp()
    app.run()