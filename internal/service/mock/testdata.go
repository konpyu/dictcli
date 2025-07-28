package mock

import "github.com/yourusername/dictcli/internal/types"

// getSampleSentences returns predefined sentences organized by topic
func getSampleSentences() map[string][]string {
	return map[string][]string{
		types.TopicBusiness: {
			"The meeting starts at nine o'clock sharp.",
			"We need to increase our quarterly revenue by fifteen percent.",
			"Please prepare the financial reports for tomorrow's presentation.",
			"Our team exceeded the sales targets this month.",
			"The deadline for this project is next Friday.",
			"We should schedule a follow-up meeting with the client.",
			"The marketing campaign was very successful last quarter.",
			"Can you send me the budget proposal by tomorrow?",
			"We're launching a new product line next year.",
			"The conference call is scheduled for three PM.",
			"Our competitor's prices are significantly lower than ours.",
			"We need to hire more software engineers for the development team.",
			"The contract negotiations will continue next week.",
			"Please review the quarterly performance metrics carefully.",
			"Our customer satisfaction scores have improved dramatically.",
		},
		types.TopicTravel: {
			"Where is the nearest subway station from here?",
			"I'd like to book a room for three nights.",
			"What time does the museum open on Sundays?",
			"How much does a taxi to the airport cost?",
			"The flight was delayed due to bad weather.",
			"We're planning to visit Tokyo next summer vacation.",
			"The hotel room has a beautiful view of the ocean.",
			"Can you recommend a good restaurant nearby?",
			"I need to exchange some money at the bank.",
			"The tour guide was very knowledgeable and friendly.",
			"We missed our connecting flight in Chicago yesterday.",
			"The passport control line was extremely long today.",
			"I packed too many clothes for this short trip.",
			"The local food was absolutely delicious and authentic.",
			"We should check the weather forecast before leaving.",
		},
		types.TopicDaily: {
			"I usually wake up at six thirty every morning.",
			"She's cooking dinner in the kitchen right now.",
			"We need to buy groceries after work today.",
			"The weather forecast says it will rain tomorrow.",
			"My younger brother is studying at university.",
			"I enjoy reading books before going to bed.",
			"We're planning a barbecue party this weekend.",
			"She takes the bus to work every single day.",
			"The children are playing in the park outside.",
			"I forgot to set my alarm clock last night.",
			"We watched a really interesting movie yesterday evening.",
			"She's learning how to drive a car finally.",
			"I need to do laundry and clean the house.",
			"We're thinking about adopting a puppy soon.",
			"The coffee shop opens at seven in the morning.",
		},
		types.TopicTechnology: {
			"The software update will be released next month.",
			"We're developing a new mobile application for Android.",
			"The server crashed due to high traffic volume.",
			"Please backup your important data regularly and consistently.",
			"The artificial intelligence system learns from user behavior.",
			"We need to upgrade our computer hardware soon.",
			"The database contains millions of customer records.",
			"Cloud storage is more secure than local storage.",
			"The video conference quality was poor today.",
			"We're implementing blockchain technology for security purposes.",
			"The new smartphone has excellent camera features.",
			"Machine learning algorithms are becoming more sophisticated.",
			"We should encrypt all sensitive data immediately.",
			"The internet connection is very slow this morning.",
			"Virtual reality is transforming the gaming industry completely.",
		},
		types.TopicHealth: {
			"You should exercise regularly to stay healthy.",
			"The doctor recommended taking vitamins every day.",
			"She's been feeling under the weather lately.",
			"We need to schedule your annual medical checkup.",
			"Eating vegetables is important for good nutrition.",
			"He broke his leg while playing football.",
			"The hospital is located downtown near the station.",
			"You should drink at least eight glasses of water daily.",
			"The dentist appointment is scheduled for next Tuesday.",
			"She's allergic to peanuts and dairy products.",
			"Regular sleep is essential for mental health.",
			"The medicine should be taken with food.",
			"We're organizing a health awareness campaign this month.",
			"Swimming is an excellent form of cardiovascular exercise.",
			"The blood test results will be ready tomorrow.",
		},
	}
}

// getGradeTemplates returns templates for generating realistic grading feedback
func getGradeTemplates() []GradeTemplate {
	return []GradeTemplate{
		{
			WER:   0.0,
			Score: 100,
			JapaneseExplanation: "完璧です！全ての単語を正確に聞き取れています。この調子で続けましょう。",
			AlternativeExpressions: []string{
				"Excellent work! Perfect dictation.",
				"Outstanding! You got every word right.",
			},
		},
		{
			WER:   0.1,
			Score: 90,
			JapaneseExplanation: "とても良い結果です。わずかなミスはありますが、全体的によく聞き取れています。",
			AlternativeExpressions: []string{
				"Great job! Just a small mistake here and there.",
				"Very good! Almost perfect dictation.",
			},
		},
		{
			WER:   0.2,
			Score: 80,
			JapaneseExplanation: "良い成果です。いくつかの単語で聞き間違いがありました。類似音に注意して練習を続けましょう。",
			AlternativeExpressions: []string{
				"Good work! Some minor errors with similar-sounding words.",
				"Well done! Keep practicing those tricky words.",
			},
		},
		{
			WER:   0.3,
			Score: 70,
			JapaneseExplanation: "まずまずの結果です。冠詞や前置詞の聞き取りに課題があるようです。短い文章から練習し直しましょう。",
			AlternativeExpressions: []string{
				"Decent attempt! Focus on articles and prepositions.",
				"Not bad! Work on those small function words.",
			},
		},
		{
			WER:   0.4,
			Score: 60,
			JapaneseExplanation: "もう少し頑張りましょう。基本的な単語は聞き取れていますが、文の構造を意識して聞くことが大切です。",
			AlternativeExpressions: []string{
				"Keep trying! You got the main words but missed some details.",
				"Practice more! Focus on sentence structure.",
			},
		},
		{
			WER:   0.5,
			Score: 50,
			JapaneseExplanation: "まだ改善の余地があります。音の変化や連結に慣れる必要があります。ゆっくりとした音声から始めてみましょう。",
			AlternativeExpressions: []string{
				"Room for improvement! Try slower audio first.",
				"Don't give up! Practice with clearer pronunciation.",
			},
		},
		{
			WER:   1.0,
			Score: 0,
			JapaneseExplanation: "今回は難しかったようですね。心配しないでください。繰り返し練習することで必ず上達します。簡単なレベルから始めてみましょう。",
			AlternativeExpressions: []string{
				"That was challenging! Don't worry, practice makes perfect.",
				"Tough one! Let's try an easier level first.",
			},
		},
	}
}