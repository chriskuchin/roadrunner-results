

## Architecture

- Races:
    - Participants: Not required but can be used for pre race registration
    - Events:
        - Results: Used for recording results and retrieving results
        - Timers: Every event has timers unless they are field.
        - Heats: Heats allow lane numbers to be associated with bib_numbers for faster result recording during track events
        - Participants:




## Workflows

- Registration
- Timer
    - Start Timer Remote
        - POST /timers Create and Start Timer
        - PUT /timers Start Existing Timer
    - Record Finish Time
        - PUT /results with timestamp of button press
- Heat/Timer Results
    - Record Bib Number/Lane Number
        - If Lane Number lookup bib number and associate it with time
- Heat
    - Record Bib and Lane Numbers
- Results
    -