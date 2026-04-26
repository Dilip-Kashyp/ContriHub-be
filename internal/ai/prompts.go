package ai

import (
	"fmt"
	"strings"
)

// BuildExplainRepoPrompt creates a prompt to explain a repository.
func BuildExplainRepoPrompt(repoName, description, language, topics, readmeSnippet, userQuestion string) string {
	prompt := fmt.Sprintf(`Explain this GitHub repository in simple terms.

**Repository:** %s
**Description:** %s
**Language:** %s
**Topics:** %s

**README (truncated):**
%s

Provide your response with these sections:
1. **What it does** — A clear 2-3 sentence explanation
2. **Tech Stack** — Languages, frameworks, and tools used
3. **Who it's for** — Target audience
4. **How to run** — Quick-start steps (if available from README)`, repoName, description, language, topics, truncateText(readmeSnippet, 1500))

	if userQuestion != "" {
		prompt += fmt.Sprintf("\n5. **Your Question** — Answer this: %s", userQuestion)
	}

	return prompt
}

// BuildFindProjectsPrompt creates a prompt to explain why projects match a query.
func BuildFindProjectsPrompt(query string, repoSummaries string) string {
	return fmt.Sprintf(`You are Gibo, an expert open-source mentor. A developer is looking for: "%s"

I have found these top repositories on GitHub for them:
%s

Analyze these projects and provide a helpful response:
1. **The Top Pick** — Identify the single best repository from this list that matches their intent. Explain exactly why it's the winner and what specific skills they can gain by contributing there.
2. **Alternative Options** — Briefly mention 2-3 other repositories from the list, explaining how they differ (e.g., "more complex architecture", "better for UI/UX", "pure backend").
3. **Getting Started Tip** — Provide one actionable piece of advice for a developer starting with their first contribution to these types of projects.

Keep the tone encouraging, technical but accessible, and highly relevant to their specific query. Output in clean markdown.`, query, truncateText(repoSummaries, 2000))
}

// BuildRoadmapPrompt creates a prompt for generating a learning roadmap.
func BuildRoadmapPrompt(interest, skillLevel string, repoSummaries string) string {
	if skillLevel == "" {
		skillLevel = "beginner"
	}
	return fmt.Sprintf(`Create a beginner-friendly learning roadmap for a %s developer interested in: **%s**

Use these real GitHub projects as learning resources:
%s

Generate an ordered learning path with:
1. **Step number and title** for each stage
2. **Which repo to study** and WHY
3. **What to learn** from each repo (specific skills/concepts)
4. **Estimated time** to spend on each step

Order from easiest to hardest. Keep it practical and actionable.`, skillLevel, interest, truncateText(repoSummaries, 2500))
}

// BuildStartGuidePrompt creates a prompt for navigating a repository.
func BuildStartGuidePrompt(repoName, description, language, readmeSnippet, fileStructure string) string {
	return fmt.Sprintf(`Help a developer get started with this repository.

**Repository:** %s
**Description:** %s  
**Language:** %s

**README (truncated):**
%s

**File Structure:**
%s

Provide a clear guide with:
1. **Entry Point** — The main file(s) to look at first
2. **Important Folders** — What each key directory contains
3. **Steps to Run** — How to set up and run the project locally
4. **Suggested Learning Order** — Which files/folders to read in what order to understand the codebase`, repoName, description, language, truncateText(readmeSnippet, 1500), truncateText(fileStructure, 800))
}

// BuildReadmePrompt creates a prompt for generating a GitHub profile README.
func BuildReadmePrompt(username, name, bio string, topRepos, languages string) string {
	return fmt.Sprintf(`Generate a professional GitHub profile README in markdown for this developer.

**Username:** %s
**Name:** %s
**Bio:** %s

**Top Repositories:**
%s

**Languages Used:**
%s

Create a polished README with these sections:
1. **Intro** — A welcoming header with name and a short tagline
2. **About Me** — 2-3 sentences based on their bio and repos
3. **Skills & Tools** — Based on their languages and repo topics (use badges/shields if appropriate)
4. **Featured Projects** — Highlight top repos with links and descriptions
5. **Contact** — GitHub profile link

Use emojis tastefully. Make it visually appealing with proper markdown formatting. Output ONLY the raw markdown, no explanations.`, username, name, bio, truncateText(topRepos, 1500), languages)
}

// BuildSummaryPrompt creates a prompt for generating a professional summary.
func BuildSummaryPrompt(skills, projects, experience string) string {
	prompt := fmt.Sprintf(`Write a concise professional summary (4-6 lines) for a software developer's resume.

**Skills:** %s
**Projects:** %s`, skills, truncateText(projects, 1000))

	if experience != "" {
		prompt += fmt.Sprintf("\n**Experience:** %s", experience)
	}

	prompt += `

The summary should:
- Be written in first person or third person (professional tone)
- Highlight key technical strengths
- Mention notable projects naturally
- Be resume-ready (no markdown formatting, plain text)
- Be 4-6 lines maximum`

	return prompt
}

// truncateText limits text to approximately maxChars characters.
func truncateText(text string, maxChars int) string {
	text = strings.TrimSpace(text)
	if len(text) <= maxChars {
		return text
	}
	return text[:maxChars] + "\n... (truncated)"
}
