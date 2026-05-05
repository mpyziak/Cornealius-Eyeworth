# Cornealius-Eyeworth
Cornealious Eyeworth hails from a lineage of distinguished butlers, each with their own unique quirks. Born in the dimly lit corridors of the Eyeworth Manor, Cornealious was destined for service. Yet, unlike his predecessors who polished silverware and arranged flowers, Cornealious found solace in the intricate world of eyes.

He has a pocket watch with a miniature iris painted on the face. When Cornealious listens, he tilts his head ever so slightly, like a raven contemplating a puzzle. His right eyebrow arches imperceptibly-an unspoken question. Is it too much screen time troubling your eyes? He knows, even before you speak. And he will help you with that.

## Cornealius' advice
###Introduction:
In the grand ballroom of life, where pixels pirouette and screens waltz, our eyes play the lead role. As your devoted ocular butler, I present to you the 20-20-20 Rule-a dance of rejuvenation for your visual senses.

### The Choreography:
Every 20 Minutes: Imagine a pocket watch ticking in the corner of your screen. When its hands align at the 20-minute mark, rise from your chair with grace. Your eyes deserve an intermission.
Step One: Look Afar: Gaze beyond the confines of your digital stage. Find a distant point-a tree, a distant building, or perhaps the horizon. Let your eyes stretch their legs.
Step Two: Blink Ballet: Blink thrice, like a prima ballerina fluttering her lashes. Blinking moisturizes your corneas, preventing them from feeling like parched parchment.
Step Three: 20 Seconds: Now, my dear friend, focus on that distant point for precisely 20 seconds. Let your retinas breathe. Imagine them sighing contentedly, like a well-steeped chamomile tea.
Step Four: Repeat Encore: Return to your screen, refreshed. The pixels will applaud your diligence. Your eyes, like seasoned performers, will bow gracefully.

### Posture and Distance: A Choreography of Proximity
Maintain good posture when viewing screens. Not too close, not too far-like a waltz partner, find the perfect distance.
Blink, my dear friend. Blink as if your eyes are whispering secrets to the universe.

Properly Adjusted Lighting: Theatrics for Screens
View screens with properly adjusted lighting. Imagine a theater stage-no harsh spotlights, just gentle illumination.

### The Ballad of Blue Light:
Ah, blue light-the rogue troubadour of our digital age. It serenades our screens, but its nocturnal melodies disrupt our slumber. Fear not, for Cornealious has a remedy:

Blue Light Filters: Adorn your screens with virtual opera glasses. Apps and settings can filter out the disruptive blue notes, allowing your eyes to rest.
Twilight Sonata: As twilight descends, dim your screens. Let them harmonize with the fading sun. Your circadian rhythm will thank you.

### Finale: The Ocular Waltz:
Picture this: You, dear reader, standing by the open window, moonlight caressing your face. Your eyes, like crystal chandeliers, sparkle with vitality. You’ve mastered the 20-20-20 Rule, and your vision sings arias of gratitude.

## Build & Publish: Summoning Cornealious from Source

Should you wish to conjure Cornealious yourself, rather than trust the pre-assembled gentleman, the following ceremony is required:

### Prerequisites: The Butler's Wardrobe
Before the performance may begin, ensure the .NET SDK is present on your machine - Cornealious will not step in without it.

### Act One: The Build
To assemble Cornealious in his casual, Debug attire:

    dotnet build -c Debug

To dress him in his finest Release garments:

    dotnet build -c Release

### Act Two: The Publish
To publish a framework-dependent build - lean and elegant, expecting the .NET runtime to already be present on the host machine:

    dotnet publish -c Release -r win-x64 -o publish --self-contained false

To publish Cornealious as a fully self-contained gentleman, runtime and all, using the provided publish profile:

    dotnet publish -c Release -r win-x64 -p:PublishProfile=win-x64-standalone

For those who prefer the Visual Studio wing of the manor: the Publish... option in the IDE will do the honours, guided by the profile found at `Properties/PublishProfiles/win-x64-standalone.pubxml`.

### Finale: A Word from Cornealious
Once published, the assembled executable and his companions will await in `publish/` (or under `bin/Release/net10.0-windows10.0.17763.0/win-x64/`). Do remember to keep `config.json` beside the executable - without it, Cornealious will stand in the hallway, hat in hand, utterly uncertain of when to advise you on eye care.
