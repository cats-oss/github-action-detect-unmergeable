'use strict';

const assert = require('assert');
const fs = require('fs').promises;

const octokit = require('@octokit/rest')();

const {
    getOpenPullRequestAll,
    checkAndMarkIfPullRequestUnmergeable,
    isRelatedToPushedBranch,
    getPullRequestId,
} = require('./operations');

(async function main() {
    const GITHUB_TOKEN = process.env.GITHUB_TOKEN;
    assertTypeof(GITHUB_TOKEN, 'string', 'GITHUB_TOKEN should be string');

    const GITHUB_REPOSITORY = process.env.GITHUB_REPOSITORY;
    assertTypeof(GITHUB_REPOSITORY, 'string', 'GITHUB_REPOSITORY should be string');

    const GITHUB_EVENT_NAME = process.env.GITHUB_EVENT_NAME;
    assert.strictEqual(GITHUB_EVENT_NAME, 'push', `${GITHUB_EVENT_NAME} event is not supported`);

    const GITHUB_EVENT_PATH = process.env.GITHUB_EVENT_PATH;
    assertTypeof(GITHUB_EVENT_PATH, 'string', 'GITHUB_EVENT_PATH should be string');

    const [REPO_OWNER, REPO_NAME] = GITHUB_REPOSITORY.split('/');
    assertTypeof(REPO_OWNER, 'string', `we could not get REPO_OWNER from ${GITHUB_REPOSITORY}`);
    assertTypeof(REPO_NAME, 'string', `we could not get REPO_NAME from ${GITHUB_REPOSITORY}`);

    octokit.authenticate({
        type: 'token',
        token: GITHUB_TOKEN
    });

    const eventDataString = await fs.readFile(GITHUB_EVENT_PATH, {
        encoding: 'utf8',
        flag: 'r'
    });
    const eventData = JSON.parse(eventDataString);

    const eventOriginRefName = eventData.ref;

    const compareUrl = eventData.compare;

    const openedPRList = await getOpenPullRequestAll(octokit, REPO_OWNER, REPO_NAME);
    const queue = [];
    for (const pullReqInfo of openedPRList) {
        const number = getPullRequestId(pullReqInfo);
        if (!isRelatedToPushedBranch(pullReqInfo, eventOriginRefName)) {
            console.log(`#${number} is not related to ${eventOriginRefName}`);
            continue;
        }

        const task = checkAndMarkIfPullRequestUnmergeable(octokit, REPO_OWNER, REPO_NAME, pullReqInfo, compareUrl);
        queue.push(task);
    }

    await queue;
})().catch((e) => {
    console.error(e);
    process.exit(1);
});

function assertTypeof(val, typename, message) {
    assert.strictEqual(typeof val, typename, message);
}