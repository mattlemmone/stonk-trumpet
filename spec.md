We're going to create a Go application that retrieves tweets from an endpoint, then passes these tweets to OpenAI for sentiment analysis related to potential stock market impact. The application will:

1. Poll for tweets only during Eastern time zone hours (7am to midnight)
2. Analyze tweet sentiment using OpenAI
3. Classify tweets as positive or negative based on potential stock market impact
4. Send notifications or text messages for tweets with significant positive impact
5. Implement a mechanism to track processed tweets using unique identifiers or timestamps
6. Ensure processed tweet tracking persists across application restarts
7. Design as a continuous-running server application

The system will prevent duplicate processing of tweets and provide real-time market sentiment insights through automated analysis and notification.
