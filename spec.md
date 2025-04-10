We're going to create a Python application that retrieves posts ("statuses") from Truth Social, then passes these posts to OpenAI for sentiment analysis related to potential stock market impact. The application will:

1. Poll for Truth Social statuses only during Eastern time zone hours (configurable, defaults to 7am to 11pm)
2. Analyze status sentiment using OpenAI's API
3. Classify statuses as positive or negative based on potential stock market impact
4. Send notifications through ntfy.sh for statuses with significant positive or negative impact
5. Implement a mechanism to track processed statuses using unique identifiers
6. Ensure processed status tracking persists across application restarts
7. Design as a continuous-running server application

The system will prevent duplicate processing of statuses and provide real-time market sentiment insights through automated analysis and notification. The ntfy.sh service will be used for push notifications when significant market impacts are detected.
