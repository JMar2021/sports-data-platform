const gamesContainer = document.querySelector("#games-container");
const sportButtons = document.querySelectorAll(".sport-button");

let games = [];

async function loadGames() {
    try {
        const response = await fetch("data/games.json");

        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`);
        }

        const data = await response.json();
        games = data.games;

        renderGames(games);
    } catch (error) {
        console.error("Failed to load games:", error);

        gamesContainer.innerHTML = `
            <p>Unable to load games.</p>
        `;
    }
}

function formatDateHeading(date) {
    const today = new Date();

    const todayStart = new Date(
        today.getFullYear(),
        today.getMonth(),
        today.getDate()
    );

    const dateStart = new Date(
        date.getFullYear(),
        date.getMonth(),
        date.getDate()
    );

    const difference =
        (dateStart - todayStart) / (1000 * 60 * 60 * 24);

    if (difference === -1) {
        return "Yesterday";
    }

    if (difference === 0) {
        return "Today";
    }

    if (difference === 1) {
        return "Tomorrow";
    }

    return date.toLocaleDateString(undefined, {
        month: "long",
        day: "numeric",
        year: "numeric",
    });
}

function formatGameTime(scheduledAt) {
    const date = new Date(scheduledAt);

    return date.toLocaleTimeString(undefined, {
        hour: "numeric",
        minute: "2-digit",
    });
}

function groupGamesByDate(gamesToGroup) {
    const groups = new Map();

    gamesToGroup.forEach((game) => {
        const date = new Date(game.scheduled_at);

        const dateKey = [
            date.getFullYear(),
            String(date.getMonth() + 1).padStart(2, "0"),
            String(date.getDate()).padStart(2, "0"),
        ].join("-");

        if (!groups.has(dateKey)) {
            groups.set(dateKey, {
                date,
                games: [],
            });
        }

        groups.get(dateKey).games.push(game);
    });

    return groups;
}

function createGameCard(game) {
    const gameCard = document.createElement("article");

    gameCard.classList.add("game-card");
    gameCard.dataset.sport = game.sport;

    const gameTime = formatGameTime(game.scheduled_at);

    gameCard.innerHTML = `
        <div class="sport-label">${game.sport.toUpperCase()}</div>

        <div class="teams">
            <div class="team">
                <span>${game.away_team.name}</span>
                <strong>${game.away_score}</strong>
            </div>

            <div class="team">
                <span>${game.home_team.name}</span>
                <strong>${game.home_score}</strong>
            </div>
        </div>

        <div class="game-info">
            <span>${game.status} · ${gameTime}</span>
            <span>${game.venue}</span>
        </div>
    `;

    return gameCard;
}

function renderGames(gamesToRender) {
    gamesContainer.innerHTML = "";

    if (gamesToRender.length === 0) {
        gamesContainer.innerHTML = "<p>No games found.</p>";
        return;
    }

    const groupedGames = groupGamesByDate(gamesToRender);

    for (const group of groupedGames.values()) {
        const dateSection = document.createElement("section");
        dateSection.classList.add("date-group");

        const dateHeading = document.createElement("h3");
        dateHeading.textContent = formatDateHeading(group.date);

        dateSection.appendChild(dateHeading);

        const gamesGrid = document.createElement("div");
        gamesGrid.classList.add("games-container");

        group.games.forEach((game) => {
            gamesGrid.appendChild(createGameCard(game));
        });

        dateSection.appendChild(gamesGrid);
        gamesContainer.appendChild(dateSection);
    }
}

sportButtons.forEach((button) => {
    button.addEventListener("click", () => {
        const selectedSport = button.dataset.sport;

        sportButtons.forEach((button) => {
            button.classList.remove("active");
        });

        button.classList.add("active");

        if (selectedSport === "all") {
            renderGames(games);
            return;
        }

        const filteredGames = games.filter(
            (game) => game.sport === selectedSport
        );

        renderGames(filteredGames);
    });
});

loadGames();